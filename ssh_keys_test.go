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

	// Create a new SSH Key
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
	defer deleteSSHKey(t, c, k.Data.ID, projectID)

	if k.Data.Attributes.Name != keyName {
		t.Fatalf("Expected new SSH key name to be %s, not %s", keyName, k.Data.Attributes.Name)
	}

	// Update newly created SSH key
	/*keyName = randString8()
	skur := SSHKeyUpdateRequest{
		Data: SSHKeyUpdateData{
			ID:   k.Data.ID,
			Type: testProjectType,
			Attributes: SSHKeyUpdateAttributes{
				Name: keyName,
			},
		},
	}
	k, _, err = c.SSHKeys.Update(k.Data.ID, projectID, &skur)
	if err != nil {
		t.Fatal(err)
	}
	if k.Data.Attributes.Name != keyName {
		t.Fatalf("Expected the name of the updated SSH key to be %s, not %s", keyName, k.Data.Attributes.Name)
	}*/

	kl, _, err := c.SSHKeys.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(kl) != 1 {
		t.Fatalf("Device List should contain exactly one device, was: %v", kl)
	}
}
