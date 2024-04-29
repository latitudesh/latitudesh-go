package latitude

import (
	"testing"
)

func TestAccAccountsBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := setup(t)
	defer stopRecord()

	t.Run("Get User Account", func(t *testing.T) {
		profile, _, err := c.Users.Get(nil)
		if err != nil {
			t.Fatal(err)
		}

		if profile == nil {
			t.Fatal("Could not find user account")
		}
	})

	t.Run("List user teams", func(t *testing.T) {
		teams, _, err := c.Users.List(nil)
		if err != nil {
			t.Fatal(err)
		}

		if len(teams) < 1 {
			t.Fatal("Team must have at least a owner")
		}
	})

}
