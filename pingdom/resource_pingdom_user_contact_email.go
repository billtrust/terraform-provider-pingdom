package pingdom

import (
	"fmt"
	"log"
	"strconv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/billtrust/go-pingdom/pingdom"
)

func resourcePingdomUserContactEmailEmail() *schema.Resource {
	return &schema.Resource{
		Create: resourcePingdomUserContactEmailCreate,
		Read:   resourcePingdomUserContactEmailRead,
		Update: resourcePingdomUserContactEmailUpdate,
		Delete: resourcePingdomUserContactEmailDelete,
		Exists: resourcePingdomUserContactEmailExists,
		Importer: &schema.ResourceImporter{
			State: resourcePingdomUserContactEmailImporter,
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

			"address": {
				Type: schema.TypeString,
				Required: true,
			},
		},
	}
}


func resourcePingdomUserContactEmailCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	contact := pingdom.Contact{}
	var userId int

	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(int)
	}

	if v, ok := d.GetOk("address"); ok {
		contact.Email = v.(string)
	}

	if v, ok := d.GetOk("severity"); ok {
		contact.Severity = v.(string)
	}

	c, err := client.Users.CreateContact(userId, contact)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(c.Id))

	return nil
}

func resourcePingdomUserContactEmailRead(d *schema.ResourceData, meta interface{}) error {
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
		for _, email := range user.Email {
			if email.Id == id {
				d.SetId(strconv.Itoa(email.Id))
				d.Set("user_id", userId)
				d.Set("address", email.Address)
				d.Set("severity", email.Severity)

				return nil
			}
		}
	}

	d.SetId("")
	return nil
}

func resourcePingdomUserContactEmailUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pingdom.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving id for resource: %s", err)
	}

	contact := pingdom.Contact{}
	var userId int

	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(int)
	}

	if v, ok := d.GetOk("address"); ok {
		contact.Email = v.(string)
	}

	if v, ok := d.GetOk("severity"); ok {
		contact.Severity = v.(string)
	}

	_, e := client.Users.UpdateContact(userId, id, contact)
	if e != nil {
		return fmt.Errorf("Error updating contact: %s", e)
	}

	return nil
}

func resourcePingdomUserContactEmailDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourcePingdomUserContactEmailExists(d *schema.ResourceData, meta interface{}) (bool, error) {
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
		for _, email := range user.Email {
			if email.Id == id {
				return true, nil
				break
			}
		}
	}
	return false, nil
}

func resourcePingdomUserContactEmailImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

	return []*schema.ResourceData{d}, nil
}

