package members_test

import (
	"testing"

	latitude "github.com/latitudesh/latitudesh-go"
	team_members "github.com/latitudesh/latitudesh-go/infrastructure/members"
)

const ()

func deleteMember(t *testing.T, c *latitude.Client, id string) {
	if _, err := c.Members.Delete(id); err != nil {
		t.Fatal(err)
	}
}

func TestAccMembersBasic(t *testing.T) {
	latitude.SkipUnlessAcceptanceTestsAllowed(t)

	c, stopRecord := latitude.Setup(t)
	defer stopRecord()
	defer latitude.ProjectTeardown(c)

	// List Members
	members, _, err := c.Members.List(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(members) < 1 {
		t.Fatal("Team must have at least a owner")
	}

	//Create Member
	request := team_members.MemberCreateRequest{
		Data: team_members.MemberCreateData{
			Type: "memberships",
			Attributes: team_members.MemberCreateAttributes{
				FirstName: "go-sdk",
				LastName:  "test",
				Email:     "go_sdk_test@latitude.sh",
				Role:      team_members.Collaborator,
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
