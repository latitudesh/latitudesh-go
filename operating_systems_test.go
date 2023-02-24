package latitude

import (
	"testing"
)

func TestAccOperatingSystemBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := setup(t)
	defer stopRecord()

	osl, _, err := c.OperatingSystems.List(nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(osl) == 0 {
		t.Fatalf("Operating System List should contain at least one plan")
	}

	OsTest := OperatingSystem{
		ID:       "6",
		Type:     "operating_system",
		Name:     "CentOS 8",
		Distro:   "centos",
		Slug:     "centos_8_x64",
		Version:  "8",
		User:     "cloud-user",
		Raid:     true,
		Rescue:   true,
		SshKeys:  true,
		UserData: true,
	}

	// Check Operating System data
	if OsTest.ID != osl[0].ID {
		t.Fatalf("Expected the id of the Operating System to be %s, not %s", OsTest.ID, osl[0].ID)
	}
	if OsTest.Type != osl[0].Type {
		t.Fatalf("Expected the type of the Operating System to be %s, not %s", OsTest.Type, osl[0].Type)
	}
	if OsTest.Name != osl[0].Name {
		t.Fatalf("Expected the name of the Operating System to be %s, not %s", OsTest.Name, osl[0].Name)
	}
}
