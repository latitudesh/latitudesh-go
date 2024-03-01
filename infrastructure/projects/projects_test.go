package projects_test

import (
	"fmt"
	"testing"

	latitude "github.com/latitudesh/latitudesh-go"
	projects "github.com/latitudesh/latitudesh-go/infrastructure/projects"
)

func TestAccProjectBasic(t *testing.T) {
	latitude.SkipUnlessAcceptanceTestsAllowed(t)

	c, stopRecord := latitude.Setup(t)
	defer stopRecord()
	defer latitude.ProjectTeardown(c)

	// List Projects
	projs, _, err := c.Projects.List(nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, proj := range projs {
		fmt.Println(proj.ID, proj.Name)
	}

	// Create a new project
	rs := latitude.TestProjectPrefix + latitude.RandString8()
	pcr := projects.ProjectCreateRequest{
		Data: projects.ProjectCreateData{
			Type: latitude.TestProjectType,
			Attributes: projects.ProjectCreateAttributes{
				Name:        rs,
				Environment: latitude.TestProjectEnvironment,
			},
		},
	}
	p, _, err := c.Projects.Create(&pcr)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != rs {
		t.Fatalf("Expected new project name to be %s, not %s", rs, p.Name)
	}

	// Update newly created project
	rs = latitude.TestProjectPrefix + latitude.RandString8()
	pur := projects.ProjectUpdateRequest{
		Data: projects.ProjectUpdateData{
			ID:   p.ID,
			Type: latitude.TestProjectType,
			Attributes: projects.ProjectCreateAttributes{
				Name:        rs,
				Environment: latitude.TestProjectEnvironment,
			},
		},
	}
	p, _, err = c.Projects.Update(p.ID, &pur)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != rs {
		t.Fatalf("Expected the name of the updated project to be %s, not %s", rs, p.Name)
	}

	// Get newly updated project
	gotProject, _, err := c.Projects.Get(p.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if gotProject.Name != rs {
		t.Fatalf("Expected the name of the GOT project to be %s, not %s", rs, gotProject.Name)
	}

	// Delete newly created project
	_, err = c.Projects.Delete(p.ID)
	if err != nil {
		t.Fatal(err)
	}
}
