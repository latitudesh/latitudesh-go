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

	var udID string

	t.Run("Create UserData", func(t *testing.T) {
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

		ud, _, err := c.UserData.Create(projectID, &udcr)
		if err != nil {
			t.Fatal(err)
		}

		udID = ud.ID

		if ud.Content != content {
			t.Fatalf("Expected new User Data content to be %s, not %s", content, ud.Content)
		}
	})

	defer deleteUserData(t, c, udID, projectID)

	t.Run("Update UserData", func(t *testing.T) {
		// Update newly created User Data
		description := randString8()
		skur := UserDataUpdateRequest{
			Data: UserDataUpdateData{
				ID:   udID,
				Type: testUserDataType,
				Attributes: UserDataUpdateAttributes{
					Description: description,
				},
			},
		}
		ud, _, err := c.UserData.Update(udID, projectID, &skur)
		if err != nil {
			t.Fatal(err)
		}
		if ud.Description != description {
			t.Fatalf("Expected the description of the updated User Data to be %s, not %s", description, ud.Description)
		}
	})

	t.Run("Get and List UserData", func(t *testing.T) {
		udl, _, err := c.UserData.List(projectID, nil)
		if err != nil {
			t.Fatal(err)
		}

		if len(udl) != 1 {
			t.Fatalf("User Data List should contain exactly one key, was: %v", udl)
		}

		ud, _, err := c.UserData.Get(udID, projectID, nil)
		if err != nil {
			t.Fatal(err)
		}
		if ud.ID != udl[0].ID {
			t.Fatalf("Expected User Data ID to be %s, not %s", udl[0].ID, ud.ID)
		}
	})
}
