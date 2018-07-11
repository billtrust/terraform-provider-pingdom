package pingdom

import (
"fmt"
"log"
"strconv"
"github.com/hashicorp/terraform/helper/schema"
"github.com/billtrust/go-pingdom/pingdom"
)

func resourcePingdomUser() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingdomUserCreate,
		Read:   resourcePingdomUserRead,
		Update: resourcePingdomUserUpdate,
		Delete: resourcePingdomUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourcePingdomUserImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"paused": {
				Type:     schema.TypeString,
				Required: false,
				ForceNew: false,
			},

			"primary": {
				Type:     schema.TypeString,
				Required: false,
				ForceNew: false,
			},

			"sms": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Resource{
					Schema: map[string]*schema.Schema{
						"severity": {
							Type:     schema.TypeString,
							Required: false,
						},
						"country_code": {
							Type: schema.TypeString,
							Required: true,
						},
						"number" : {
							Type: schema.TypeString,
							Required: true,
						},
						"provider" : {
							Type: schema.TypeString,
							Required: false,
						},
					},
				},
			},
			"email": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem:     &schema.Resource{
					Schema: map[string]*schema.Schema{
						"severity": {
							Type:     schema.TypeString,
							Required: false,
						},
						"address": {
							Type: schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func checkForUserResource(d *schema.ResourceData) (pingdom.User, error) {
	userParams := pingdom.User{}

	// required
	if v, ok := d.GetOk("name"); ok {
		userParams.Username = v.(string)
	}

	if userParams.Username == "" {
		return pingdom.User{}, fmt.Errorf("name cannot be blank")
	}

	if v, ok := d.GetOk("primary"); ok {
		userParams.Primary = v.(string)
	}

	if v, ok := d.GetOk("paused"); ok {
		userParams.Paused = v.(string)
	}

	if v, ok := d.GetOk("sms"); ok {
		smsList := v.(*schema.Set).List()
		var smsSlice []pingdom.UserSmsResponse
		for i := range smsList {
			smsSlice = append(smsSlice, smsList[i].(pingdom.UserSmsResponse))
		}
		userParams.Sms = smsSlice
	}

	if v, ok := d.GetOk("email"); ok {
		emailList := v.(*schema.Set).List()
		var emailSlice []pingdom.UserEmailResponse
		for i := range emailList {
			emailSlice = append(emailSlice, emailList[i].(pingdom.UserEmailResponse))
		}
		userParams.Email = emailSlice
	}

	return userParams, nil
}


func resourcePingdomUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	user, err := checkForUserResource(d)
	if err != nil {
		return err
	}

	u, err := client.Users.Create(&user)
	if err != nil {
		return err
	}

	for i:= range user.Sms {
		contact := pingdom.Contact{
			Provider : user.Sms[i].Provider,
			Number : user.Sms[i].Number,
			CountryCode : user.Sms[i].CountryCode,
			Severity : user.Sms[i].Severity,
		}
		client.Users.CreateContact(u.Id, contact)
	}

	for i:= range user.Email {
		contact := pingdom.Contact{
			Severity : user.Email[i].Severity,
			Email : user.Email[i].Address,
		}
		client.Users.CreateContact(u.Id, contact)
	}


	d.SetId(strconv.Itoa(u.Id))

	return nil
}

func resourcePingdomUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	users, err := client.Users.List()
	if err != nil {
		return fmt.Errorf("Error retrieving list of users: %s", err)
	}

	var foundUser pingdom.UsersResponse
	exists := false
	for _, user := range users {
		if user.Id == id {
			foundUser = user
			exists = true
			break
		}
	}
	if !exists {
		d.SetId("")
		return nil
	}

	d.Set("name", foundUser.Username)
	d.Set("paused", foundUser.Paused)

	smsSet := []interface{}{}
	for _, sms := range foundUser.Sms {
		item := map[string]interface{}{
			"severity" : sms.Severity,
			"country_code" : sms.CountryCode,
			"number" : sms.Number,
			"provider" : sms.Provider,
		}
		smsSet = append(smsSet, item)
	}

	emailSet := []interface{}{}
	for _, email := range foundUser.Email {
		item := map[string]interface{}{
			"severity" : email.Severity,
			"address" : email.Address,
		}
		emailSet = append(emailSet, item)
	}

	return nil
}

func resourcePingdomUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	user, err := checkForUserResource(d)
	if err != nil {
		return err
	}

	_, err = client.Users.Update(id, &user)
	if err != nil {
		return fmt.Errorf("Error updating user: %s", err)
	}

	for i:= range user.Sms {
		contact := pingdom.Contact{
			Provider : user.Sms[i].Provider,
			Number : user.Sms[i].Number,
			CountryCode : user.Sms[i].CountryCode,
			Severity : user.Sms[i].Severity,
		}
		client.Users.UpdateContact(id, user.Sms[i].Id, contact)
	}

	for i:= range user.Email {
		contact := pingdom.Contact{
			Severity : user.Email[i].Severity,
			Email : user.Email[i].Address,
		}
		client.Users.UpdateContact(id, user.Email[i].Id, contact)
	}


	return nil
}

func resourcePingdomUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	log.Printf("[INFO] Deleting User: %v", id)

	_, err = client.Users.Delete(id)
	if err != nil {
		return fmt.Errorf("Error deleting user: %s", err)
	}

	return nil
}

func resourcePingdomUserImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	return []*schema.ResourceData{d}, nil
}

