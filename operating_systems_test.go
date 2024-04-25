package latitude

import (
	"testing"
)

func TestAccOperatingSystemBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	c, stopRecord := setup(t)
	defer stopRecord()

	t.Run("List Operating Systems", func(t *testing.T) {
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

			assertEqual(t, os.Type, osTest.Type, "Operating System Type")
			assertEqual(t, os.Name, osTest.Name, "Operating System Name")
			assertEqual(t, os.Distro, osTest.Distro, "Operating System Distro")
			assertEqual(t, os.Slug, osTest.Slug, "Operating System Slug")
			assertEqual(t, os.Version, osTest.Version, "Operating System Version")
			assertEqual(t, os.User, osTest.User, "Operating System User")
			assertEqual(t, os.Raid, osTest.Raid, "Operating System Raid")
			assertEqual(t, os.Rescue, osTest.Rescue, "Operating System Rescue")
			assertEqual(t, os.SshKeys, osTest.SshKeys, "Operating System SshKeys")
			assertEqual(t, os.UserData, osTest.UserData, "Operating System UserData")

			return
		}
		t.Fatalf("Operating System with id %s not found", osTest.ID)
	})
}
