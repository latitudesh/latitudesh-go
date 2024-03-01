package projects

import (
	"path"

	api "github.com/latitudesh/latitudesh-go/api_utils"
	internal "github.com/latitudesh/latitudesh-go/internal"
	types "github.com/latitudesh/latitudesh-go/types"
)

const ProjectBasePath = "/projects"

// ProjectService interface defines available project methods
type ProjectService interface {
	List(listOpt *api.ListOptions) ([]Project, *types.Response, error)
	Get(string, *api.GetOptions) (*Project, *types.Response, error)
	Create(*ProjectCreateRequest) (*Project, *types.Response, error)
	Update(string, *ProjectUpdateRequest) (*Project, *types.Response, error)
	Delete(string) (*types.Response, error)
}

type ProjectRoot struct {
	Data ProjectData   `json:"data"`
	Meta internal.Meta `json:"meta"`
}

type ProjectData struct {
	ID         string               `json:"id"`
	Type       string               `json:"type"`
	Attributes ProjectGetAttributes `json:"attributes"`
}

type ProjectListResponse struct {
	Data []ProjectData `json:"data"`
	Meta internal.Meta `json:"meta"`
}

type ProjectGetResponse struct {
	Data ProjectData   `json:"data"`
	Meta internal.Meta `json:"meta"`
}

type ProjectGetAttributes struct {
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Description   string `json:"description"`
	BillingType   string `json:"billing_type"`
	BillingMethod string `json:"billing_method"`
	Environment   string `json:"environment"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// ProjectCreateRequest type used to create a Latitude project
type ProjectCreateRequest struct {
	Data ProjectCreateData `json:"data"`
}

type ProjectCreateData struct {
	Type       string                  `json:"type"`
	Attributes ProjectCreateAttributes `json:"attributes"`
}

type ProjectCreateAttributes struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Environment string `json:"environment"`
}

// ProjectUpdateRequest type used to update a Latitude project
type ProjectUpdateRequest struct {
	Data ProjectUpdateData `json:"data"`
}

type ProjectUpdateData struct {
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Attributes ProjectCreateAttributes `json:"attributes"`
}

// ProjectServiceOp implements ProjectService
type ProjectServiceOp struct {
	Client internal.RequestDoer
}

type Project struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Description   string `json:"description"`
	BillingType   string `json:"billing_type"`
	BillingMethod string `json:"billing_method"`
	Environment   string `json:"environment"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// Flatten latitude API data structures
func NewFlatProject(pd ProjectData) Project {
	return Project{
		pd.ID,
		pd.Attributes.Name,
		pd.Attributes.Slug,
		pd.Attributes.Description,
		pd.Attributes.BillingType,
		pd.Attributes.BillingMethod,
		pd.Attributes.Environment,
		pd.Attributes.CreatedAt,
		pd.Attributes.UpdatedAt,
	}
}

func NewFlatProjectList(pd []ProjectData) []Project {
	var res []Project
	for _, project := range pd {
		res = append(res, NewFlatProject(project))
	}
	return res
}

// List returns a list of projects
func (s *ProjectServiceOp) List(opts *api.ListOptions) (projects []Project, resp *types.Response, err error) {
	apiPathQuery := opts.WithQuery(ProjectBasePath)

	for {
		res := new(ProjectListResponse)

		resp, err = s.Client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		projects = append(projects, NewFlatProjectList(res.Data)...)

		if apiPathQuery = api.NextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

// Get returns a project by id
func (s *ProjectServiceOp) Get(projectID string, opts *api.GetOptions) (*Project, *types.Response, error) {
	endpointPath := path.Join(ProjectBasePath, projectID)
	apiPathQuery := opts.WithQuery(endpointPath)
	project := new(ProjectGetResponse)
	resp, err := s.Client.DoRequest("GET", apiPathQuery, nil, project)
	if err != nil {
		return nil, resp, err
	}

	flatProject := NewFlatProject(project.Data)
	return &flatProject, resp, err
}

// Create creates a new project
func (s *ProjectServiceOp) Create(createRequest *ProjectCreateRequest) (*Project, *types.Response, error) {
	project := new(ProjectGetResponse)

	resp, err := s.Client.DoRequest("POST", ProjectBasePath, createRequest, project)
	if err != nil {
		return nil, resp, err
	}

	flatProject := NewFlatProject(project.Data)
	return &flatProject, resp, err
}

// Update updates a project
func (s *ProjectServiceOp) Update(projectID string, updateRequest *ProjectUpdateRequest) (*Project, *types.Response, error) {
	apiPath := path.Join(ProjectBasePath, projectID)
	project := new(ProjectGetResponse)

	resp, err := s.Client.DoRequest("PATCH", apiPath, updateRequest, project)
	if err != nil {
		return nil, resp, err
	}

	flatProject := NewFlatProject(project.Data)
	return &flatProject, resp, err
}

// Delete deletes a project
func (s *ProjectServiceOp) Delete(projectID string) (*types.Response, error) {
	apiPath := path.Join(ProjectBasePath, projectID)

	return s.Client.DoRequest("DELETE", apiPath, nil, nil)
}
