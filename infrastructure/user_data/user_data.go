package user_data

import (
	"path"

	api "github.com/latitudesh/latitudesh-go/api_utils"
	internal "github.com/latitudesh/latitudesh-go/internal"
	types "github.com/latitudesh/latitudesh-go/types"

	projects "github.com/latitudesh/latitudesh-go/infrastructure/projects"
)

const userDataBasePath = "/user_data"

type UserDataService interface {
	List(projectID string, opts *api.ListOptions) ([]UserData, *types.Response, error)
	Get(userDataID, projectID string, opts *api.GetOptions) (*UserData, *types.Response, error)
	Create(projectID string, request *UserDataCreateRequest) (*UserData, *types.Response, error)
	Update(userDataID, projectID string, request *UserDataUpdateRequest) (*UserData, *types.Response, error)
	Delete(userDataID, projectID string) (*types.Response, error)
}

// UserData represents a Latitude User Data record
type UserData struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Content     string `json:"content"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// UserDataServiceOp implements UserDataService
type UserDataServiceOp struct {
	Client internal.RequestDoer
}

// UserDataCreateRequest type used to create a Latitude User Data record
type UserDataCreateRequest struct {
	Data UserDataCreateData `json:"data"`
}

type UserDataCreateData struct {
	Type       string                   `json:"type"`
	Attributes UserDataCreateAttributes `json:"attributes"`
}

type UserDataCreateAttributes struct {
	Description string `json:"description"`
	Content     string `json:"content"`
}

type UserDataUpdateRequest struct {
	Data UserDataUpdateData `json:"data"`
}

type UserDataUpdateData struct {
	ID         string                   `json:"id"`
	Type       string                   `json:"type"`
	Attributes UserDataUpdateAttributes `json:"attributes"`
}

type UserDataUpdateAttributes struct {
	Description string `json:"description,omitempty"`
	Content     string `json:"content,omitempty"`
}

type UserDataGetResponse struct {
	Data UserDataData  `json:"data"`
	Meta internal.Meta `json:"meta"`
}

type UserDataData struct {
	ID         string                `json:"id"`
	Type       string                `json:"type"`
	Attributes UserDataGetAttributes `json:"attributes"`
}

type UserDataGetAttributes struct {
	Description string `json:"description"`
	Content     string `json:"content"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UserDataListResponse struct {
	Data []UserDataData `json:"data"`
	Meta internal.Meta  `json:"meta"`
}

// Flatten latitude API data structures
func NewFlatUserData(ud UserDataData) UserData {
	return UserData{
		ud.ID,
		ud.Attributes.Description,
		ud.Attributes.Content,
		ud.Attributes.CreatedAt,
		ud.Attributes.UpdatedAt,
	}
}

func NewFlatUserDataList(udList []UserDataData) []UserData {
	var userDataList []UserData

	for _, userData := range udList {
		userDataList = append(userDataList, NewFlatUserData(userData))
	}

	return userDataList
}

// List returns list of User data
func (u *UserDataServiceOp) List(projectID string, opts *api.ListOptions) ([]UserData, *types.Response, error) {
	var userDataList []UserData
	endpointPath := path.Join(projects.ProjectBasePath, projectID, userDataBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		userDataRecords := new(UserDataListResponse)

		resp, err := u.Client.DoRequest("GET", apiPathQuery, nil, userDataRecords)
		if err != nil {
			return nil, resp, err
		}

		userDataList = append(userDataList, NewFlatUserDataList(userDataRecords.Data)...)

		if apiPathQuery = api.NextPage(userDataRecords.Meta, opts); apiPathQuery != "" {
			continue
		}

		return userDataList, resp, err
	}

}

// Get returns a User data by id
func (u *UserDataServiceOp) Get(userDataID, projectID string, opts *api.ListOptions) (*UserData, *types.Response, error) {
	endpointPath := path.Join(projects.ProjectBasePath, projectID, userDataBasePath, userDataID)
	apiPathQuery := opts.WithQuery(endpointPath)
	userData := new(UserDataGetResponse)

	resp, err := u.Client.DoRequest("GET", apiPathQuery, nil, userData)
	if err != nil {
		return nil, resp, err
	}

	flatUserData := NewFlatUserData(userData.Data)
	return &flatUserData, resp, err
}

// Create creates a new User Data record
func (s *UserDataServiceOp) Create(projectID string, createRequest *UserDataCreateRequest) (*UserData, *types.Response, error) {
	endpointPath := path.Join(projects.ProjectBasePath, projectID, userDataBasePath)
	userData := new(UserDataGetResponse)

	resp, err := s.Client.DoRequest("POST", endpointPath, createRequest, userData)
	if err != nil {
		return nil, resp, err
	}

	flatUserData := NewFlatUserData(userData.Data)
	return &flatUserData, resp, err
}

// Update updates a User Data record
func (s *UserDataServiceOp) Update(userDataID, projectID string, updateRequest *UserDataUpdateRequest) (*UserData, *types.Response, error) {
	apiPath := path.Join(projects.ProjectBasePath, projectID, userDataBasePath, userDataID)
	userData := new(UserDataGetResponse)

	resp, err := s.Client.DoRequest("PATCH", apiPath, updateRequest, userData)
	if err != nil {
		return nil, resp, err
	}

	flatUserData := NewFlatUserData(userData.Data)
	return &flatUserData, resp, err
}

// Delete deletes a User Data record
func (s *UserDataServiceOp) Delete(userDataID, projectID string) (*types.Response, error) {
	apiPath := path.Join(projects.ProjectBasePath, projectID, userDataBasePath, userDataID)

	return s.Client.DoRequest("DELETE", apiPath, nil, nil)
}
