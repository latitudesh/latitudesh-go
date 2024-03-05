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

	t.Run("Create Team", func(t *testing.T) {
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

		team, _, err := c.Teams.Create(&tcr)
		if err != nil {
			t.Fatal(err)
		}

		if team.Description != description {
			t.Fatalf("Expected team description to be %s, not %s", description, team.Description)
		}
	})

	t.Run("Get Team", func(t *testing.T) {
        _, _, err := c.Teams.Get()
		if err != nil {
			t.Fatal(err)
		}
	})
}
