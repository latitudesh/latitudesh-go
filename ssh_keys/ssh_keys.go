package ssh_keys

import (
	"path"

	api "github.com/latitudesh/latitudesh-go/api_utils"
	internal "github.com/latitudesh/latitudesh-go/internal"
	types "github.com/latitudesh/latitudesh-go/types"

	projects "github.com/latitudesh/latitudesh-go/projects"
	servers "github.com/latitudesh/latitudesh-go/servers"
)

const sshKeyBasePath = "/ssh_keys"

type SSHKeyService interface {
	List(projectID string, opts *api.ListOptions) ([]SSHKey, *types.Response, error)
	Get(sshKeyID string, projectID string, opts *api.GetOptions) (*SSHKey, *types.Response, error)
	Create(projectID string, request *SSHKeyCreateRequest) (*SSHKey, *types.Response, error)
	Update(sshKeyID string, projectID string, request *SSHKeyUpdateRequest) (*SSHKey, *types.Response, error)
	Delete(sshKeyID string, projectID string) (*types.Response, error)
}

type SSHKeyRoot struct {
	Data servers.ServerData `json:"data"`
	Meta internal.Meta      `json:"meta"`
}

type SSHKeyData struct {
	ID         string              `json:"id"`
	Type       string              `json:"type"`
	Attributes SSHKeyGetAttributes `json:"attributes"`
}

type SSHKeyListResponse struct {
	Data []SSHKeyData  `json:"data"`
	Meta internal.Meta `json:"meta"`
}

type SSHKeyGetResponse struct {
	Data SSHKeyData    `json:"data"`
	Meta internal.Meta `json:"meta"`
}

type SSHKeyGetAttributes struct {
	Name        string `json:"name"`
	PublicKey   string `json:"public_key"`
	Fingerprint string `json:"fingerprint"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
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
	Name string `json:"name"`
}

// SSHKeyServiceOp implements SSHKeyService
type SSHKeyServiceOp struct {
	Client internal.RequestDoer
}

// SSHKey represents a Latitude SSH key
type SSHKey struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PublicKey   string `json:"public_key"`
	Fingerprint string `json:"fingerprint"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
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
func (s *SSHKeyServiceOp) List(projectID string, opts *api.ListOptions) (sshKeys []SSHKey, resp *types.Response, err error) {
	endpointPath := path.Join(projects.ProjectBasePath, projectID, sshKeyBasePath)
	apiPathQuery := opts.WithQuery(endpointPath)

	for {
		res := new(SSHKeyListResponse)

		resp, err = s.Client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}
		sshKeys = append(sshKeys, NewFlatSSHKeyList(res.Data)...)

		if apiPathQuery = api.NextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

// Get returns an SSH key by id
func (s *SSHKeyServiceOp) Get(sshKeyID string, projectID string, opts *api.GetOptions) (*SSHKey, *types.Response, error) {
	endpointPath := path.Join(projects.ProjectBasePath, projectID, sshKeyBasePath, sshKeyID)
	apiPathQuery := opts.WithQuery(endpointPath)
	sshKey := new(SSHKeyGetResponse)
	resp, err := s.Client.DoRequest("GET", apiPathQuery, nil, sshKey)
	if err != nil {
		return nil, resp, err
	}

	flatSSHKey := NewFlatSSHKey(sshKey.Data)
	return &flatSSHKey, resp, err
}

// Create creates a new SSH key
func (s *SSHKeyServiceOp) Create(projectID string, createRequest *SSHKeyCreateRequest) (*SSHKey, *types.Response, error) {
	endpointPath := path.Join(projects.ProjectBasePath, projectID, sshKeyBasePath)
	sshKey := new(SSHKeyGetResponse)

	resp, err := s.Client.DoRequest("POST", endpointPath, createRequest, sshKey)
	if err != nil {
		return nil, resp, err
	}

	flatSSHKey := NewFlatSSHKey(sshKey.Data)
	return &flatSSHKey, resp, err
}

// Update updates an SSH key
func (s *SSHKeyServiceOp) Update(sshKeyID string, projectID string, updateRequest *SSHKeyUpdateRequest) (*SSHKey, *types.Response, error) {
	apiPath := path.Join(projects.ProjectBasePath, projectID, sshKeyBasePath, sshKeyID)
	sshKey := new(SSHKeyGetResponse)

	resp, err := s.Client.DoRequest("PATCH", apiPath, updateRequest, sshKey)
	if err != nil {
		return nil, resp, err
	}

	flatSSHKey := NewFlatSSHKey(sshKey.Data)
	return &flatSSHKey, resp, err
}

// Delete deletes an SSH Key
func (s *SSHKeyServiceOp) Delete(sshKeyID string, projectID string) (*types.Response, error) {
	apiPath := path.Join(projects.ProjectBasePath, projectID, sshKeyBasePath, sshKeyID)

	return s.Client.DoRequest("DELETE", apiPath, nil, nil)
}
