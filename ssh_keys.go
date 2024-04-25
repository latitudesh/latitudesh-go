package latitude

import (
	"path"
)

const sshKeyBasePath = "/ssh_keys"

type SSHKeyService interface {
	List(projectID string, opts *ListOptions) ([]SSHKey, *Response, error)
	Get(sshKeyID string, projectID string, opts *GetOptions) (*SSHKey, *Response, error)
	Create(projectID string, request *SSHKeyCreateRequest) (*SSHKey, *Response, error)
	Update(sshKeyID string, projectID string, request *SSHKeyUpdateRequest) (*SSHKey, *Response, error)
	Delete(sshKeyID string, projectID string) (*Response, error)
}

type SSHKeyRoot struct {
	Data ServerData `json:"data"`
	Meta meta       `json:"meta"`
}

type SSHKeyData struct {
	ID         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes SSHKeyGetAttributes `json:"attributes"`
}

type SSHKeyListResponse struct {
	Data []SSHKeyData `json:"data"`
	Meta meta         `json:"meta"`
}

type SSHKeyGetResponse struct {
	Data SSHKeyData `json:"data"`
	Meta meta       `json:"meta"`
}

type SSHKeyGetAttributes struct {
	Name        string     `json:"name"`
	PublicKey   string     `json:"public_key"`
	Fingerprint string     `json:"fingerprint"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	Tags        []EmbedTag `json:"tags"`
}

// SSHKeyCreateRequest type used to create a Latitude SSH key
type SSHKeyCreateRequest struct {
	Data SSHKeyCreateData `json:"data"`
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
	Name string   `json:"name"`
	Tags []string `json:"tags,omitempty"`
}

// SSHKeyServiceOp implements SSHKeyService
type SSHKeyServiceOp struct {
	client requestDoer
}

// SSHKey represents a Latitude SSH key
type SSHKey struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	PublicKey   string     `json:"public_key"`
	Fingerprint string     `json:"fingerprint"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	Tags        []EmbedTag `json:"tags"`
}

// Flatten latitude API data structures
func NewFlatSSHKey(sd SSHKeyData) SSHKey {
	return SSHKey{
		sd.ID,
		sd.Attributes.Name,
		sd.Attributes.PublicKey,
		sd.Attributes.Fingerprint,
		sd.Attributes.CreatedAt,
		sd.Attributes.UpdatedAt,
		sd.Attributes.Tags,
	}
}

func NewFlatSSHKeyList(sd []SSHKeyData) []SSHKey {
	var res []SSHKey
	for _, ssh_key := range sd {
		res = append(res, NewFlatSSHKey(ssh_key))
	}
	return res
}

// List returns a list of SSH Keys
func (s *SSHKeyServiceOp) List(projectID string, opts *ListOptions) (sshKeys []SSHKey, resp *Response, err error) {
	endpointPath := path.Join(projectBasePath, projectID, sshKeyBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		res := new(SSHKeyListResponse)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}
		sshKeys = append(sshKeys, NewFlatSSHKeyList(res.Data)...)

		if apiPathQuery = nextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

// Get returns an SSH key by id
func (s *SSHKeyServiceOp) Get(sshKeyID string, projectID string, opts *GetOptions) (*SSHKey, *Response, error) {
	endpointPath := path.Join(projectBasePath, projectID, sshKeyBasePath, sshKeyID)
	apiPathQuery := opts.WithQuery(endpointPath)
	sshKey := new(SSHKeyGetResponse)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, sshKey)
	if err != nil {
		return nil, resp, err
	}

	flatSSHKey := NewFlatSSHKey(sshKey.Data)
	return &flatSSHKey, resp, err
}

// Create creates a new SSH key
func (s *SSHKeyServiceOp) Create(projectID string, createRequest *SSHKeyCreateRequest) (*SSHKey, *Response, error) {
	endpointPath := path.Join(projectBasePath, projectID, sshKeyBasePath)
	sshKey := new(SSHKeyGetResponse)

	resp, err := s.client.DoRequest("POST", endpointPath, createRequest, sshKey)
	if err != nil {
		return nil, resp, err
	}

	flatSSHKey := NewFlatSSHKey(sshKey.Data)
	return &flatSSHKey, resp, err
}

// Update updates an SSH key
func (s *SSHKeyServiceOp) Update(sshKeyID string, projectID string, updateRequest *SSHKeyUpdateRequest) (*SSHKey, *Response, error) {
	apiPath := path.Join(projectBasePath, projectID, sshKeyBasePath, sshKeyID)
	sshKey := new(SSHKeyGetResponse)

	resp, err := s.client.DoRequest("PATCH", apiPath, updateRequest, sshKey)
	if err != nil {
		return nil, resp, err
	}

	flatSSHKey := NewFlatSSHKey(sshKey.Data)
	return &flatSSHKey, resp, err
}

// Delete deletes an SSH Key
func (s *SSHKeyServiceOp) Delete(sshKeyID string, projectID string) (*Response, error) {
	apiPath := path.Join(projectBasePath, projectID, sshKeyBasePath, sshKeyID)

	return s.client.DoRequest("DELETE", apiPath, nil, nil)
}
