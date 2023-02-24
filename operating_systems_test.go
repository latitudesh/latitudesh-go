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
		t.Fatalf("Operating System List should contain at least one OS")
	}

	OsTest := OperatingSystem{
		ID:       "7",
		Type:     "operating_system",
		Name:     "CentOS 7.4",
		Distro:   "centos",
		Slug:     "centos_7_4_x64",
		Version:  "7.4",
		User:     "centos",
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
	if OsTest.Distro != osl[0].Distro {
		t.Fatalf("Expected the name of the Operating System to be %s, not %s", OsTest.Distro, osl[0].Distro)
	}
	if OsTest.Slug != osl[0].Slug {
		t.Fatalf("Expected the name of the Operating System to be %s, not %s", OsTest.Slug, osl[0].Slug)
	}
	if OsTest.Version != osl[0].Version {
		t.Fatalf("Expected the name of the Operating System to be %s, not %s", OsTest.Version, osl[0].Version)
	}
	if OsTest.User != osl[0].User {
		t.Fatalf("Expected the name of the Operating System to be %s, not %s", OsTest.User, osl[0].User)
	}
	if OsTest.Raid != osl[0].Raid {
		t.Fatalf("Expected the name of the Operating System to be %t, not %t", OsTest.Raid, osl[0].Raid)
	}
	if OsTest.Rescue != osl[0].Rescue {
		t.Fatalf("Expected the name of the Operating System to be %t, not %t", OsTest.Rescue, osl[0].Rescue)
	}
	if OsTest.SshKeys != osl[0].SshKeys {
		t.Fatalf("Expected the name of the Operating System to be %t, not %t", OsTest.SshKeys, osl[0].SshKeys)
	}
	if OsTest.UserData != osl[0].UserData {
		t.Fatalf("Expected the name of the Operating System to be %t, not %t", OsTest.UserData, osl[0].UserData)
	}
}
