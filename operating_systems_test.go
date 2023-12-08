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

	osTest := OperatingSystem{
		ID:       "os_KgXQvNe3azpbP",
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
	for _, os := range osl {
		if os.ID != osTest.ID {
			continue
		}

		if osTest.Type != os.Type {
			t.Fatalf("Expected the type of the Operating System to be %s, not %s", osTest.Type, os.Type)
		}
		if osTest.Name != os.Name {
			t.Fatalf("Expected the name of the Operating System to be %s, not %s", osTest.Name, os.Name)
		}
		if osTest.Distro != os.Distro {
			t.Fatalf("Expected the name of the Operating System to be %s, not %s", osTest.Distro, os.Distro)
		}
		if osTest.Slug != os.Slug {
			t.Fatalf("Expected the name of the Operating System to be %s, not %s", osTest.Slug, os.Slug)
		}
		if osTest.Version != os.Version {
			t.Fatalf("Expected the name of the Operating System to be %s, not %s", osTest.Version, os.Version)
		}
		if osTest.User != os.User {
			t.Fatalf("Expected the name of the Operating System to be %s, not %s", osTest.User, os.User)
		}
		if osTest.Raid != os.Raid {
			t.Fatalf("Expected the name of the Operating System to be %t, not %t", osTest.Raid, os.Raid)
		}
		if osTest.Rescue != os.Rescue {
			t.Fatalf("Expected the name of the Operating System to be %t, not %t", osTest.Rescue, os.Rescue)
		}
		if osTest.SshKeys != os.SshKeys {
			t.Fatalf("Expected the name of the Operating System to be %t, not %t", osTest.SshKeys, os.SshKeys)
		}
		if osTest.UserData != os.UserData {
			t.Fatalf("Expected the name of the Operating System to be %t, not %t", osTest.UserData, os.UserData)
		}
		return
	}
	t.Fatalf("Operating System with id %s not found", osTest.ID)
}
