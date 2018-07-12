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
		Exists: resourcePingdomUserExists,
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
				Optional: true,
				ForceNew: false,
			},

			"primary": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
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

	d.SetId(strconv.Itoa(u.Id))

	return nil
}

func resourcePingdomUserExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return false, fmt.Errorf("Error retrieving id for resource: %s", err)
	}
	users, err := client.Users.List()
	if err != nil {
		return false, fmt.Errorf("Error retrieving list of users: %s", err)
	}

	exists := false
	for _, user := range users {
		if user.Id == id {
			exists = true
			break
		}
	}
	return exists, nil
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

