package latitude

import (
	"testing"
)

func TestAccAccountsBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := setup(t)
	defer stopRecord()

	//Get User
	_, _, err := c.Users.Get(nil)
	if err != nil {
		t.Fatal(err)
	}

}
