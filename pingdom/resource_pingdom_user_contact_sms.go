package pingdom

import (
	"fmt"
	"log"
	"strconv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/billtrust/go-pingdom/pingdom"
)

func resourcePingdomUserContactSms() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingdomUserContactSmsCreate,
		Read:   resourcePingdomUserContactSmsRead,
		Update: resourcePingdomUserContactSmsUpdate,
		Delete: resourcePingdomUserContactSmsDelete,
		Exists: resourcePingdomUserContactSmsExists,
		Importer: &schema.ResourceImporter{
			State: resourcePingdomUserContactSmsImporter,
		},

		Schema: map[string]*schema.Schema{
			"user_id" : {
				Type: schema.TypeInt,
				Required: true,
			},
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
				Type:     schema.TypeString,
				Required: false,
			},
		},
	}
}


func resourcePingdomUserContactSmsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	contact := pingdom.Contact{}
	var userId int

	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(int)
	}

	if v, ok := d.GetOk("provider"); ok {
		contact.Provider = v.(string)
	}

	if v, ok := d.GetOk("number"); ok {
		contact.Number = v.(string)
	}

	if v, ok := d.GetOk("country_code"); ok {
		contact.CountryCode = v.(string)
	}

	if v, ok := d.GetOk("severity"); ok {
		contact.Severity = v.(string)
	}

	resp, err := client.Users.CreateContact(userId, contact)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(resp.Id))

	return nil
}

func resourcePingdomUserContactSmsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	users, err := client.Users.List()
	if err != nil {
		return fmt.Errorf("Error retrieving list of users: %s", err)
	}

	for _, user := range users {
		userId := user.Id
		for _, sms := range user.Sms {
			if sms.Id == id {
				d.SetId(strconv.Itoa(sms.Id))
				d.Set("user_id", userId)
				d.Set("provider", sms.Provider)
				d.Set("number", sms.Number)
				d.Set("country_code", sms.CountryCode)
				d.Set("severity", sms.Severity)

				return nil
			}
		}
	}

	d.SetId("")
	return nil
}

func resourcePingdomUserContactSmsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)
	var userId int
	contact := pingdom.Contact{}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(int)
	}

	if v, ok := d.GetOk("provider"); ok {
		contact.Provider = v.(string)
	}

	if v, ok := d.GetOk("number"); ok {
		contact.Number = v.(string)
	}

	if v, ok := d.GetOk("country_code"); ok {
		contact.CountryCode = v.(string)
	}

	if v, ok := d.GetOk("severity"); ok {
		contact.Severity = v.(string)
	}

	log.Printf("[INFO] Updating User Contact: %v", id)

	_, err = client.Users.UpdateContact(userId, id, contact)
	if err != nil {
		return fmt.Errorf("Error updating contact: %s", err)
	}

	return nil
}

func resourcePingdomUserContactSmsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	var userId int
	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(int)
	}

	log.Printf("[INFO] Deleting User Contact: %v", id)

	_, err = client.Users.DeleteContact(userId, id)
	if err != nil {
		return fmt.Errorf("Error deleting user contact: %s", err)
	}

	return nil
}

func resourcePingdomUserContactSmsExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	users, err := client.Users.List()
	if err != nil {
		return false, fmt.Errorf("Error retrieving list of users and contacts: %s", err)
	}

	for _, user := range users {
		for _, sms := range user.Sms {
			if sms.Id == id {
				return true, nil
				break
			}
		}
	}
	return false, nil
}

func resourcePingdomUserContactSmsImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	return []*schema.ResourceData{d}, nil
}

