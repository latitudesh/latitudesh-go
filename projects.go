package latitude

import (
	"path"
)

const projectBasePath = "/projects"

// ProjectService interface defines available project methods
type ProjectService interface {
	List(listOpt *ListOptions) ([]ProjectData, *Response, error)
	Get(string, *GetOptions) (*Project, *Response, error)
	Create(*ProjectCreateRequest) (*Project, *Response, error)
	Update(string, *ProjectUpdateRequest) (*Project, *Response, error)
	Delete(string) (*Response, error)
}

// Project represents a Latitude project
type Project struct {
	Data ProjectData `json:"data"`
	Meta meta        `json:"meta"`
}

type ProjectData struct {
	ID         string               `json:"id"`
	Type       string            `json:"type"`
	Attributes ProjectAttributes `json:"attributes"`
}

type ProjectAttributes struct {
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Description   string `json:"description"`
	BillingType   string `json:"billing_type"`
	BillingMethod string `json:"billing_method"`
	Environment   string `json:"environment"`
}

type ProjectListResponse struct {
	Projects []ProjectData `json:"data"`
	Meta     meta          `json:"meta"`
}

// ProjectCreateRequest type used to create a Latitude project
type ProjectCreateRequest struct {
	Data ProjectCreateData `json:"data"`
}

func (p ProjectCreateRequest) String() string {
	return Stringify(p)
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
	ID         string                     `json:"id"`
	Type       string                  `json:"type"`
	Attributes ProjectCreateAttributes `json:"attributes"`
}

func (p ProjectUpdateRequest) String() string {
	return Stringify(p)
}

// ProjectServiceOp implements ProjectService
type ProjectServiceOp struct {
	client requestDoer
}

// List returns a list of projects
func (s *ProjectServiceOp) List(opts *ListOptions) (projects []ProjectData, resp *Response, err error) {
	apiPathQuery := opts.WithQuery(projectBasePath)

	for {
		subset := new(ProjectListResponse)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		projects = append(projects, subset.Projects...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

// Get returns a project by id
func (s *ProjectServiceOp) Get(projectID string, opts *GetOptions) (*Project, *Response, error) {
	endpointPath := path.Join(projectBasePath, projectID)
	apiPathQuery := opts.WithQuery(endpointPath)
	project := new(Project)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, project)
	if err != nil {
		return nil, resp, err
	}
	return project, resp, err
}

// Create creates a new project
func (s *ProjectServiceOp) Create(createRequest *ProjectCreateRequest) (*Project, *Response, error) {
	project := new(Project)

	resp, err := s.client.DoRequest("POST", projectBasePath, createRequest, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, err
}

// Update updates a project
func (s *ProjectServiceOp) Update(projectID string, updateRequest *ProjectUpdateRequest) (*Project, *Response, error) {
	apiPath := path.Join(projectBasePath, projectID)
	project := new(Project)

	resp, err := s.client.DoRequest("PATCH", apiPath, updateRequest, project)
	if err != nil {
		return nil, resp, err
	}

	return project, resp, err
}

// Delete deletes a project
func (s *ProjectServiceOp) Delete(projectID string) (*Response, error) {
	apiPath := path.Join(projectBasePath, projectID)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}
