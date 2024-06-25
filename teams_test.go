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
		name := randString8()
		address := randString8()

		tcr := TeamCreateRequest{
			Data: TeamCreateData{
				Type: testTeamType,
				Attributes: TeamCreateAttributes{
					Name:     name,
					Currency: "USD",
					Address:  address,
				},
			},
		}

		team, _, err := c.Teams.Create(&tcr)
		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, team.Name, name, "Team Name")
		assertEqual(t, *team.Address, address, "Team Address")
	})

	t.Run("Get Team", func(t *testing.T) {
		_, _, err := c.Teams.Get()
		if err != nil {
			t.Fatal(err)
		}
	})
}
