package latitude

import (
	"testing"
)

func TestAccRolesBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := setup(t)
	defer stopRecord()

	// List Roles
	roles, _, err := c.Roles.List(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(roles) < 1 {
		t.Fatal("Team must have at least a role")
	}

	//Get Role
	role, _, err := c.Roles.Get(roles[0].ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Check Role data
	if role.ID != roles[0].ID {
		t.Fatalf("Expected the id of the GOT role to be %s, not %s", roles[0].ID, role.ID)
	}
	if role.Name != roles[0].Name {
		t.Fatalf("Expected the line of the GOT plan to be %s, not %s", roles[0].Name, role.Name)
	}
}
