package latitude

import (
	"path"
)

const UserDataBasePath = "/user_data"

type UserDataService interface {
	List(projectID string, opts *ListOptions) ([]UserDataData, *Response, error)
	Get(UserDataID string, projectID string, opts *GetOptions) (*UserDataGetResponse, *Response, error)
	Create(projectID string, request *UserDataCreateRequest) (*UserData, *Response, error)
	Update(UserDataID string, projectID string, request *UserDataUpdateRequest) (*UserData, *Response, error)
	Delete(UserDataID string, projectID string) (*Response, error)
}

type UserData struct {
	Data UserDataData `json:"data"`
	Meta meta         `json:"meta"`
}

type UserDataData struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Attributes UserDataAttributes `json:"attributes"`
}

type UserDataAttributes struct {
	Name        string `json:"name"`
	PublicKey   string `json:"public_key"`
	Fingerprint string `json:"fingerprint"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UserDataListResponse struct {
	UserData []UserDataData `json:"data"`
	Meta     meta           `json:"meta"`
}

type UserDataGetResponse struct {
	Data UserDataGetData `json:"data"`
	Meta meta            `json:"meta"`
}

type UserDataGetData struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Attributes UserDataAttributes `json:"attributes"`
}

type UserDataCreateRequest struct {
	Data UserDataCreateData `json:"data"`
}

func (s UserDataCreateRequest) String() string {
	return Stringify(s)
}

type UserDataCreateData struct {
	Type       string                   `json:"type"`
	Attributes UserDataCreateAttributes `json:"attributes"`
}

type UserDataCreateAttributes struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

// ProjectUpdateRequest type used to update a Latitude project
type UserDataUpdateRequest struct {
	Data UserDataUpdateData `json:"data"`
}

type UserDataUpdateData struct {
	ID         string                   `json:"id"`
	Type       string                   `json:"type"`
	Attributes UserDataUpdateAttributes `json:"attributes"`
}

type UserDataUpdateAttributes struct {
	Name string `json:"name"`
}

func (p UserDataUpdateRequest) String() string {
	return Stringify(p)
}

type UserDataServiceOp struct {
	client requestDoer
}

func (s *UserDataServiceOp) List(projectID string, opts *ListOptions) (userData []UserDataData, resp *Response, err error) {
	endpointPath := path.Join(projectBasePath, projectID, UserDataBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		subset := new(UserDataListResponse)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

func (s *UserDataServiceOp) Get(UserDataID string, projectID string, opts *GetOptions) (*UserDataGetResponse, *Response, error) {
	endpointPath := path.Join(projectBasePath, projectID, UserDataBasePath, UserDataID)
	apiPathQuery := opts.WithQuery(endpointPath)
	userData := new(UserDataGetResponse)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, userData)
	if err != nil {
		return nil, resp, err
	}
	return userData, resp, err
}

func (s *UserDataServiceOp) Create(projectID string, createRequest *UserDataCreateRequest) (*UserData, *Response, error) {
	endpointPath := path.Join(projectBasePath, projectID, UserDataBasePath)
	userData := new(UserData)

	resp, err := s.client.DoRequest("POST", endpointPath, createRequest, userData)
	if err != nil {
		return nil, resp, err
	}

	return userData, resp, err
}

func (s *UserDataServiceOp) Update(UserDataID string, projectID string, updateRequest *UserDataUpdateRequest) (*UserData, *Response, error) {
	apiPath := path.Join(projectBasePath, projectID, UserDataBasePath, UserDataID)
	userData := new(UserData)

	resp, err := s.client.DoRequest("PATCH", apiPath, updateRequest, userData)
	if err != nil {
		return nil, resp, err
	}

	return userData, resp, err
}

func (s *UserDataServiceOp) Delete(userDataID string, projectID string) (*Response, error) {
	apiPath := path.Join(projectBasePath, projectID, UserDataBasePath, userDataID)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}
