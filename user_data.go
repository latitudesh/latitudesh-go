package latitude

import (
	"path"
)

const userDataBasePath = "/user_data"

type UserDataService interface {
	List(projectID string, opts *ListOptions) ([]UserData, *Response, error)
	Get(userDataID, projectID string, opts *GetOptions) (*UserData, *Response, error)
	Create(projectID string, request *UserDataCreateRequest) (*UserData, *Response, error)
	Update(userDataID, projectID string, request *UserDataUpdateRequest) (*UserData, *Response, error)
	Delete(userDataID, projectID string) (*Response, error)
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
	client requestDoer
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
	Data UserDataData `json:"data"`
	Meta meta         `json:"meta"`
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
	Meta meta           `json:"meta"`
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
func (u *UserDataServiceOp) List(projectID string, opts *ListOptions) ([]UserData, *Response, error) {
	var userDataList []UserData
	endpointPath := path.Join(projectBasePath, projectID, userDataBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		userDataRecords := new(UserDataListResponse)

		resp, err := u.client.DoRequest("GET", apiPathQuery, nil, userDataRecords)
		if err != nil {
			return nil, resp, err
		}

		userDataList = append(userDataList, NewFlatUserDataList(userDataRecords.Data)...)

		if apiPathQuery = nextPage(userDataRecords.Meta, opts); apiPathQuery != "" {
			continue
		}

		return userDataList, resp, err
	}

}

// Get returns a User data by id
func (u *UserDataServiceOp) Get(userDataID, projectID string, opts *ListOptions) (*UserData, *Response, error) {
	endpointPath := path.Join(projectBasePath, projectID, userDataBasePath, userDataID)
	apiPathQuery := opts.WithQuery(endpointPath)
	userData := new(UserDataGetResponse)

	resp, err := u.client.DoRequest("GET", apiPathQuery, nil, userData)
	if err != nil {
		return nil, resp, err
	}

	flatUserData := NewFlatUserData(userData.Data)
	return &flatUserData, resp, err
}

// Create creates a new User Data record
func (s *UserDataServiceOp) Create(projectID string, createRequest *UserDataCreateRequest) (*UserData, *Response, error) {
	endpointPath := path.Join(projectBasePath, projectID, userDataBasePath)
	userData := new(UserDataGetResponse)

	resp, err := s.client.DoRequest("POST", endpointPath, createRequest, userData)
	if err != nil {
		return nil, resp, err
	}

	flatUserData := NewFlatUserData(userData.Data)
	return &flatUserData, resp, err
}

// Update updates a User Data record
func (s *UserDataServiceOp) Update(userDataID, projectID string, updateRequest *UserDataUpdateRequest) (*UserData, *Response, error) {
	apiPath := path.Join(projectBasePath, projectID, userDataBasePath, userDataID)
	userData := new(UserDataGetResponse)

	resp, err := s.client.DoRequest("PATCH", apiPath, updateRequest, userData)
	if err != nil {
		return nil, resp, err
	}

	flatUserData := NewFlatUserData(userData.Data)
	return &flatUserData, resp, err
}

// Delete deletes a User Data record
func (s *UserDataServiceOp) Delete(userDataID, projectID string) (*Response, error) {
	apiPath := path.Join(projectBasePath, projectID, userDataBasePath, userDataID)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}
