package latitude

import (
	"testing"
)

const (
	testUserDataType = "user_data"
)

func deleteUserData(t *testing.T, c *Client, userDataID string, projectID string) {
	if _, err := c.UserData.Delete(userDataID, projectID); err != nil {
		t.Fatal(err)
	}
}

func TestAccUserDataBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)
	t.Parallel()

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	// Create a new UserData record
	description := randString8()
	content := testUserDataContent()

	udcr := UserDataCreateRequest{
		Data: UserDataCreateData{
			Type: testUserDataType,
			Attributes: UserDataCreateAttributes{
				Description: description,
				Content:     content,
			},
		},
	}

	k, _, err := c.UserData.Create(projectID, &udcr)
	if err != nil {
		t.Fatal(err)
	}
	defer deleteUserData(t, c, k.ID, projectID)

	if k.Content != content {
		t.Fatalf("Expected new User Data content to be %s, not %s", content, k.Content)
	}

	// Update newly created User Data
	description = randString8()
	skur := UserDataUpdateRequest{
		Data: UserDataUpdateData{
			ID:   k.ID,
			Type: testUserDataType,
			Attributes: UserDataUpdateAttributes{
				Description: description,
			},
		},
	}
	k, _, err = c.UserData.Update(k.ID, projectID, &skur)
	if err != nil {
		t.Fatal(err)
	}
	if k.Description != description {
		t.Fatalf("Expected the description of the updated User Data to be %s, not %s", description, k.Description)
	}

	// List newly created User Data
	kl, _, err := c.UserData.List(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(kl) != 1 {
		t.Fatalf("User Data List should contain exactly one key, was: %v", kl)
	}

	// Get newly created User Data
	k, _, err = c.UserData.Get(k.ID, projectID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if k.ID != kl[0].ID {
		t.Fatalf("Expected User Data ID to be %s, not %s", kl[0].ID, k.ID)
	}
}
