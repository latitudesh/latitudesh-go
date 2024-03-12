package latitude

import (
	"testing"
)

const (
	testProjectType        = "projects"
	testProjectEnvironment = "Development"
)

func deleteProject(t *testing.T, c *Client, id string) {
	if _, err := c.Projects.Delete(id); err != nil {
		t.Fatal(err)
	}
}

func TestAccProjectBasic(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, stopRecord := setup(t)
	defer stopRecord()
	defer projectTeardown(c)

	tagIDs, deleteTags := setupTestTags(t, c)
	defer deleteTags()

	var projectID string
	t.Run("Create Project", func(t *testing.T) {
		// Create a new project
		rs := testProjectPrefix + randString8()
		pcr := ProjectCreateRequest{
			Data: ProjectCreateData{
				Type: testProjectType,
				Attributes: ProjectCreateAttributes{
					Name:        rs,
					Environment: testProjectEnvironment,
				},
			},
		}
		p, _, err := c.Projects.Create(&pcr)
		if err != nil {
			t.Fatal(err)
		}

		projectID = p.ID

		assertEqual(t, p.Name, rs, "Project Name")
	})

	defer deleteProject(t, c, projectID)

	var projectName string
	t.Run("Update Project", func(t *testing.T) {
		rs := testProjectPrefix + randString8()
		pur := ProjectUpdateRequest{
			Data: ProjectUpdateData{
				ID:   projectID,
				Type: testProjectType,
				Attributes: ProjectUpdateAttributes{
					Name:        rs,
					Environment: testProjectEnvironment,
					Tags:        tagIDs,
				},
			},
		}

		p, _, err := c.Projects.Update(projectID, &pur)
		if err != nil {
			t.Fatal(err)
		}

		projectName = p.Name
		assertEqual(t, projectName, rs, "Project Name")
		assertEqual(t, len(p.Tags), 2, "Project Tags")
	})

	t.Run("Get Project", func(t *testing.T) {
		gotProject, _, err := c.Projects.Get(projectID, nil)
		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, gotProject.Name, projectName, "Project Name")
	})

	t.Run("List Project", func(t *testing.T) {
		projs, _, err := c.Projects.List(nil)
		if err != nil {
			t.Fatal(err)
		}

		if len(projs) == 0 {
			t.Fatalf("Project List should contain at least one project")
		}
	})
}
