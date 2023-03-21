package latitude

import (
	"testing"
)

const (
	testTeamType = "teams"
)

func TestAccTeamBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, teardown := setup(t)
	defer teardown()

	// Create a new Team record
	description := randString8()
	name := randString8()
	address := randString8()

	tcr := TeamCreateRequest{
		Data: TeamCreateData{
			Type: testTeamType,
			Attributes: TeamCreateAttributes{
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
