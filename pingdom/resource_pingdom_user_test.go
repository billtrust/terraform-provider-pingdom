package pingdom

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/magiconair/properties/assert"
	"github.com/billtrust/go-pingdom/pingdom"
	"fmt"
)

func TestPingdomUser_CheckForUserResource(t *testing.T) {

	name := "testUsername"
	primary := "Y"
	paused := "N"

	rd := schema.ResourceData{}
	rd.Set("name", name)
	rd.Set("primary", primary)
	rd.Set("paused", paused)

	expectedUser := pingdom.User{
		Username: name,
		Primary: primary,
		Paused: paused,
	}

	user, err := checkForUserResource(&rd)
	assert.Equal(t, err, nil, "Error should be empty")
	assert.Equal(t, user, expectedUser, "User should be returned.")
}

func TestPingdomUser_CheckForUserResource_Fail(t *testing.T) {

	name := ""
	primary := "Y"
	paused := "N"

	rd := schema.ResourceData{}
	rd.Set("name", name)
	rd.Set("primary", primary)
	rd.Set("paused", paused)

	expectedUser := pingdom.User{}
	expectedErr := fmt.Errorf("name cannot be blank")

	user, err := checkForUserResource(&rd)
	assert.Equal(t, err, expectedErr, "Error should be empty")
	assert.Equal(t, user, expectedUser, "User should be returned.")
}
