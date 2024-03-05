package latitude

import (
	"testing"
)

const (
	testSSHKeyType = "ssh_keys"
)

func deleteSSHKey(t *testing.T, c *Client, sshKeyID string, projectID string) {
	if _, err := c.SSHKeys.Delete(sshKeyID, projectID); err != nil {
		t.Fatal(err)
	}
}

func TestAccSSHKeyBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	var keyID string

	t.Run("Create SSH Key", func(t *testing.T) {
		keyName := randString8()
		skcr := SSHKeyCreateRequest{
			Data: SSHKeyCreateData{
				Type: testSSHKeyType,
				Attributes: SSHKeyCreateAttributes{
					Name:      keyName,
					PublicKey: testSSHKey(),
				},
			},
		}

		k, _, err := c.SSHKeys.Create(projectID, &skcr)
		if err != nil {
			t.Fatal(err)
		}

		if k.Name != keyName {
			t.Fatalf("Expected new SSH key name to be %s, not %s", keyName, k.Name)
		}

		keyID = k.ID
	})

	defer deleteSSHKey(t, c, keyID, projectID)

	t.Run("Get and List SSHKeys", func(t *testing.T) {
		kList, _, err := c.SSHKeys.List(projectID, nil)
		if err != nil {
			t.Fatal(err)
		}

		if len(kList) == 0 {
			t.Fatalf("SSH key List should contain at least one key")
		}

		// Get first listed SSHkey
		gotKey, _, err := c.SSHKeys.Get(kList[0].ID, projectID, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Check SSHKey data
		if gotKey.ID != kList[0].ID {
			t.Fatalf("Expected the id of the GOT key to be %s, not %s", kList[0].ID, gotKey.ID)
		}
		if gotKey.Name != kList[0].Name {
			t.Fatalf("Expected the Name of the GOT key to be %s, not %s", kList[0].Name, gotKey.Name)
		}
		if gotKey.PublicKey != kList[0].PublicKey {
			t.Fatalf("Expected the name of the GOT key to be %s, not %s", kList[0].PublicKey, gotKey.PublicKey)
		}
	})

	t.Run("Update SSH Key", func(t *testing.T) {
		keyName := randString8()
		skur := SSHKeyUpdateRequest{
			Data: SSHKeyUpdateData{
				ID:   keyID,
				Type: testSSHKeyType,
				Attributes: SSHKeyUpdateAttributes{
					Name: keyName,
				},
			},
		}
		k, _, err := c.SSHKeys.Update(keyID, projectID, &skur)
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
	})
}
