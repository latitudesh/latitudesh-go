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

		assertEqual(t, k.Name, keyName, "SSH Key Name")

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
		assertEqual(t, gotKey.ID, kList[0].ID, "SSH Key ID")
		assertEqual(t, gotKey.Name, kList[0].Name, "SSH Key Name")
		assertEqual(t, gotKey.PublicKey, kList[0].PublicKey, "SSH Key PublicKey")
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
		assertEqual(t, k.Name, keyName, "SSH Key Name")

		kl, _, err := c.SSHKeys.List(projectID, nil)
		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, len(kl), 1, "SSH Key List length")
	})
}
