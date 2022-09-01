package latitude

import (
	"path"
)

const sshKeyBasePath = "/ssh_keys"

type SSHKeyService interface {
	List(projectID string, opts *ListOptions) ([]SSHKeyData, *Response, error)
	Get(sshKeyID string, projectID string, opts *GetOptions) (*SSHKeyGetResponse, *Response, error)
	Create(projectID string, request *SSHKeyCreateRequest) (*SSHKey, *Response, error)
	Update(sshKeyID string, projectID string, request *SSHKeyUpdateRequest) (*SSHKey, *Response, error)
	Delete(sshKeyID string, projectID string) (*Response, error)
}

// SSHKey represents a Latitude SSH key
type SSHKey struct {
	Data SSHKeyData `json:"data"`
	Meta meta       `json:"meta"`
}

type SSHKeyData struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes SSHKeyAttributes `json:"attributes"`
}

type SSHKeyAttributes struct {
	Name        string `json:"name"`
	PublicKey   string `json:"public_key"`
	Fingerprint string `json:"fingerprint"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type SSHKeyListResponse struct {
	SSHKeys []SSHKeyData `json:"data"`
	Meta    meta         `json:"meta"`
}

type SSHKeyGetResponse struct {
	Data SSHKeyGetData `json:"data"`
	Meta meta          `json:"meta"`
}

type SSHKeyGetData struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes SSHKeyAttributes `json:"attributes"`
}

// SSHKeyCreateRequest type used to create a Latitude SSH key
type SSHKeyCreateRequest struct {
	Data SSHKeyCreateData `json:"data"`
}

func (s SSHKeyCreateRequest) String() string {
	return Stringify(s)
}

type SSHKeyCreateData struct {
	Type       string                 `json:"type"`
	Attributes SSHKeyCreateAttributes `json:"attributes"`
}

type SSHKeyCreateAttributes struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

// ProjectUpdateRequest type used to update a Latitude project
type SSHKeyUpdateRequest struct {
	Data SSHKeyUpdateData `json:"data"`
}

type SSHKeyUpdateData struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes SSHKeyUpdateAttributes `json:"attributes"`
}

type SSHKeyUpdateAttributes struct {
	Name string `json:"name"`
}

func (p SSHKeyUpdateRequest) String() string {
	return Stringify(p)
}

// SSHKeyServiceOp implements SSHKeyService
type SSHKeyServiceOp struct {
	client requestDoer
}

// List returns a list of SSH Keys
func (s *SSHKeyServiceOp) List(projectID string, opts *ListOptions) (sshKeys []SSHKeyData, resp *Response, err error) {
	endpointPath := path.Join(projectBasePath, projectID, sshKeyBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		subset := new(SSHKeyListResponse)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}
		sshKeys = append(sshKeys, subset.SSHKeys...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

// Get returns an SSH key by id
func (s *SSHKeyServiceOp) Get(sshKeyID string, projectID string, opts *GetOptions) (*SSHKeyGetResponse, *Response, error) {
	endpointPath := path.Join(projectBasePath, projectID, sshKeyBasePath, sshKeyID)
	apiPathQuery := opts.WithQuery(endpointPath)
	sshKey := new(SSHKeyGetResponse)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, sshKey)
	if err != nil {
		return nil, resp, err
	}
	return sshKey, resp, err
}

// Create creates a new SSH key
func (s *SSHKeyServiceOp) Create(projectID string, createRequest *SSHKeyCreateRequest) (*SSHKey, *Response, error) {
	endpointPath := path.Join(projectBasePath, projectID, sshKeyBasePath)
	sshKey := new(SSHKey)

	resp, err := s.client.DoRequest("POST", endpointPath, createRequest, sshKey)
	if err != nil {
		return nil, resp, err
	}

	return sshKey, resp, err
}

// Update updates an SSH key
func (s *SSHKeyServiceOp) Update(sshKeyID string, projectID string, updateRequest *SSHKeyUpdateRequest) (*SSHKey, *Response, error) {
	apiPath := path.Join(projectBasePath, projectID, sshKeyBasePath, sshKeyID)
	sshKey := new(SSHKey)

	resp, err := s.client.DoRequest("PATCH", apiPath, updateRequest, sshKey)
	if err != nil {
		return nil, resp, err
	}

	return sshKey, resp, err
}

// Delete deletes an SSH Key
func (s *SSHKeyServiceOp) Delete(sshKeyID string, projectID string) (*Response, error) {
	apiPath := path.Join(projectBasePath, projectID, sshKeyBasePath, sshKeyID)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}
