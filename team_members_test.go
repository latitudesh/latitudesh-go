package latitude

import (
	"testing"
)

const ()

func deleteMember(t *testing.T, c *Client, id string) {
	if _, err := c.Members.Delete(id); err != nil {
		t.Fatal(err)
	}
}

func TestAccMembersBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, stopRecord := setup(t)
	defer stopRecord()
	defer projectTeardown(c)

	// List Members
	members, _, err := c.Members.List(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(members) < 1 {
		t.Fatal("Team must have at least a owner")
	}

	//Create Member
	request := MemberCreateRequest{
		Data: MemberCreateData{
			Type: "memberships",
			Attributes: MemberCreateAttributes{
				FirstName: "go-sdk",
				LastName:  "test",
				Email:     "go_sdk_test@latitude.sh",
				Role:      Collaborator,
			},
		},
	}

	m, _, err := c.Members.Create(&request)
	if err != nil {
		t.Fatal(err)
	}

	//Delete Member
	deleteMember(t, c, m.ID)
}
