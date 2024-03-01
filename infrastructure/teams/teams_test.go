package teams_test

import (
	"testing"

	latitude "github.com/latitudesh/latitudesh-go"
	teams "github.com/latitudesh/latitudesh-go/infrastructure/teams"
)

const (
	testTeamType = "teams"
)

func TestAccTeamBasic(t *testing.T) {
	latitude.SkipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, teardown := latitude.Setup(t)
	defer teardown()

	// Create a new Team record
	description := latitude.RandString8()
	name := latitude.RandString8()
	address := latitude.RandString8()

	tcr := teams.TeamCreateRequest{
		Data: teams.TeamCreateData{
			Type: testTeamType,
			Attributes: teams.TeamCreateAttributes{
				Description: description,
				Name:        name,
				Currency:    "USD",
				Address:     address,
			},
		},
	}

	k, _, err := c.Teams.Create(&tcr)
	if err != nil {
		t.Fatal(err)
	}

	if k.Description != description {
		t.Fatalf("Expected team description to be %s, not %s", description, k.Description)
	}

	// Get Team record
	_, _, err = c.Teams.Get()
	if err != nil {
		t.Fatal(err)
	}
}
