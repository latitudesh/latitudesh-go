package ssh_keys_test

import (
	"testing"

	latitude "github.com/latitudesh/latitudesh-go"
	sshkeys "github.com/latitudesh/latitudesh-go/infrastructure/ssh_keys"
)

const (
	testSSHKeyType = "ssh_keys"
)

func deleteSSHKey(t *testing.T, c *latitude.Client, sshKeyID string, projectID string) {
	if _, err := c.SSHKeys.Delete(sshKeyID, projectID); err != nil {
		t.Fatal(err)
	}
}

func TestAccSSHKeyBasic(t *testing.T) {
	latitude.SkipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := latitude.SetupWithProject(t)
	defer teardown()

	// Create a new SSH Key
	keyName := latitude.RandString8()
	skcr := sshkeys.SSHKeyCreateRequest{
		Data: sshkeys.SSHKeyCreateData{
			Type: testSSHKeyType,
			Attributes: sshkeys.SSHKeyCreateAttributes{
				Name:      keyName,
				PublicKey: latitude.TestSSHKey(),
			},
		},
	}

	k, _, err := c.SSHKeys.Create(projectID, &skcr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteSSHKey(t, c, k.ID, projectID)

	if k.Name != keyName {
		t.Fatalf("Expected new SSH key name to be %s, not %s", keyName, k.Name)
	}

	kList, _, err := c.SSHKeys.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(kList) == 0 {
		t.Fatalf("Plan List should contain at least one plan")
	}

	// Get first listed SSHkey
	gotKey, _, err := c.SSHKeys.Get(kList[0].ID, projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Check SSHKey data
	if gotKey.ID != kList[0].ID {
		t.Fatalf("Expected the id of the GOT plan to be %s, not %s", kList[0].ID, gotKey.ID)
	}
	if gotKey.Name != kList[0].Name {
		t.Fatalf("Expected the Name of the GOT plan to be %s, not %s", kList[0].Name, gotKey.Name)
	}
	if gotKey.PublicKey != kList[0].PublicKey {
		t.Fatalf("Expected the name of the GOT plan to be %s, not %s", kList[0].PublicKey, gotKey.PublicKey)
	}

	// Update newly created SSH key
	keyName = latitude.RandString8()
	skur := sshkeys.SSHKeyUpdateRequest{
		Data: sshkeys.SSHKeyUpdateData{
			ID:   k.ID,
			Type: testSSHKeyType,
			Attributes: sshkeys.SSHKeyUpdateAttributes{
				Name: keyName,
			},
		},
	}
	k, _, err = c.SSHKeys.Update(k.ID, projectID, &skur)
	if err != nil {
		t.Fatal(err)
	}
	if k.Name != keyName {
		t.Fatalf("Expected the name of the updated SSH key to be %s, not %s", keyName, k.Name)
	}

	kl, _, err := c.SSHKeys.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(kl) != 1 {
		t.Fatalf("SSH key List should contain exactly one key, was: %v", kl)
	}
}
