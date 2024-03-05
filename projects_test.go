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

		if p.Name != rs {
			t.Fatalf("Expected new project name to be %s, not %s", rs, p.Name)
		}
	})

	defer deleteProject(t, c, projectID)

	var projectName string
	t.Run("Update Project", func(t *testing.T) {
		rs := testProjectPrefix + randString8()
		pur := ProjectUpdateRequest{
			Data: ProjectUpdateData{
				ID:   projectID,
				Type: testProjectType,
				Attributes: ProjectCreateAttributes{
					Name:        rs,
					Environment: testProjectEnvironment,
				},
			},
		}

		p, _, err := c.Projects.Update(projectID, &pur)
		if err != nil {
			t.Fatal(err)
		}

		projectName = p.Name

		if p.Name != rs {
			t.Fatalf("Expected the name of the updated project to be %s, not %s", rs, p.Name)
		}
	})

	t.Run("Get Project", func(t *testing.T) {
		gotProject, _, err := c.Projects.Get(projectID, nil)
		if err != nil {
			t.Fatal(err)
		}

		if gotProject.Name != projectName {
			t.Fatalf("Expected the name of the GOT project to be %s, not %s", projectName, gotProject.Name)
		}
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
