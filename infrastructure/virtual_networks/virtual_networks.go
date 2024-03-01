package virtual_networks

import (
	"path"

	api "github.com/latitudesh/latitudesh-go/api_utils"
	internal "github.com/latitudesh/latitudesh-go/internal"
	types "github.com/latitudesh/latitudesh-go/types"
)

const virtualNetworkBasePath = "/virtual_networks"

type VirtualNetworkService interface {
	List(listOpt *api.ListOptions) ([]VirtualNetwork, *types.Response, error)
	Get(virtualNetworkID string, getOpt *api.GetOptions) (*VirtualNetwork, *types.Response, error)
	Create(createRequest *VirtualNetworkCreateRequest) (*VirtualNetwork, *types.Response, error)
	Update(virtualNetworkID string, updateRequest *VirtualNetworkUpdateRequest) (*VirtualNetwork, *types.Response, error)
	Delete(virtualNetworkID string) (*types.Response, error)
}

type VirtualNetworkServiceOp struct {
	Client internal.RequestDoer
}

type VirtualNetworkData struct {
	ID         string                   `json:"id"`
	Type       string                   `json:"type"`
	Attributes VirtualNetworkAttributes `json:"attributes"`
}

type VirtualNetworkAttributes struct {
	Vid              int                  `json:"vid"`
	Description      string               `json:"description"`
	Region           VirtualNetworkRegion `json:"region"`
	AssignmentsCount int                  `json:"assignments_count"`
}

type VirtualNetworkRegion struct {
	City    string             `json:"city"`
	Country string             `json:"country"`
	Site    VirtualNetworkSite `json:"site"`
}

type VirtualNetworkSite struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Facility string `json:"facility"`
}

type VirtualNetwork struct {
	ID               string `json:"id"`
	Type             string `json:"type"`
	Vid              int    `json:"vid"`
	Description      string `json:"description"`
	City             string `json:"city"`
	Country          string `json:"country"`
	SiteId           string `json:"site_id"`
	SiteName         string `json:"site_name"`
	SiteSlug         string `json:"site_slug"`
	Facility         string `json:"facility"`
	AssignmentsCount int    `json:"assignments_count"`
}

type VirtualNetworkListResponse struct {
	Data []VirtualNetworkData `json:"data"`
	Meta internal.Meta        `json:"meta"`
}

type VirtualNetworkGetResponse struct {
	Data VirtualNetworkData `json:"data"`
	Meta internal.Meta      `json:"meta"`
}

type VirtualNetworkCreateRequest struct {
	Data VirtualNetworkCreateData `json:"data"`
}

type VirtualNetworkCreateData struct {
	Type       string                         `json:"type"`
	Attributes VirtualNetworkCreateAttributes `json:"attributes"`
}

type VirtualNetworkCreateAttributes struct {
	Description string `json:"description"`
	Site        string `json:"site"`
	Project     string `json:"project"`
}

type VirtualNetworkUpdateRequest struct {
	Data VirtualNetworkUpdateData `json:"data"`
}

type VirtualNetworkUpdateData struct {
	ID         string                         `json:"id"`
	Type       string                         `json:"type"`
	Attributes VirtualNetworkUpdateAttributes `json:"attributes"`
}

type VirtualNetworkUpdateAttributes struct {
	Description string `json:"description"`
}

func NewFlatVirtualNetwork(vnd VirtualNetworkData) VirtualNetwork {
	return VirtualNetwork{
		vnd.ID,
		vnd.Type,
		vnd.Attributes.Vid,
		vnd.Attributes.Description,
		vnd.Attributes.Region.City,
		vnd.Attributes.Region.Country,
		vnd.Attributes.Region.Site.ID,
		vnd.Attributes.Region.Site.Name,
		vnd.Attributes.Region.Site.Slug,
		vnd.Attributes.Region.Site.Facility,
		vnd.Attributes.AssignmentsCount,
	}
}

func NewFlatVirtualNetworkList(vnd []VirtualNetworkData) []VirtualNetwork {
	var res []VirtualNetwork
	for _, vn := range vnd {
		res = append(res, NewFlatVirtualNetwork(vn))
	}
	return res
}

func (vn *VirtualNetworkServiceOp) List(opts *api.ListOptions) (virtualNetworks []VirtualNetwork, resp *types.Response, err error) {
	apiPathQuery := opts.WithQuery(virtualNetworkBasePath)

	for {
		res := new(VirtualNetworkListResponse)

		resp, err = vn.Client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		virtualNetworks = append(virtualNetworks, NewFlatVirtualNetworkList(res.Data)...)

		if apiPathQuery = api.NextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}

// Get returns a server by id
func (s *VirtualNetworkServiceOp) Get(virtualNetworkID string, opts *api.GetOptions) (*VirtualNetwork, *types.Response, error) {
	endpointPath := path.Join(virtualNetworkBasePath, virtualNetworkID)
	apiPathQuery := opts.WithQuery(endpointPath)
	virtualNetwork := new(VirtualNetworkGetResponse)

	resp, err := s.Client.DoRequest("GET", apiPathQuery, nil, virtualNetwork)
	if err != nil {
		return nil, resp, err
	}

	flatVirtualNetwork := NewFlatVirtualNetwork(virtualNetwork.Data)
	return &flatVirtualNetwork, resp, err
}

// Create creates a new virtual network
func (s *VirtualNetworkServiceOp) Create(createRequest *VirtualNetworkCreateRequest) (*VirtualNetwork, *types.Response, error) {
	vLan := new(VirtualNetworkGetResponse)

	resp, err := s.Client.DoRequest("POST", virtualNetworkBasePath, createRequest, vLan)
	if err != nil {
		return nil, resp, err
	}

	flatVirtualNetwork := NewFlatVirtualNetwork(vLan.Data)
	return &flatVirtualNetwork, resp, err
}

// Update updates a virtual network
func (s *VirtualNetworkServiceOp) Update(virtualNetworkID string, updateRequest *VirtualNetworkUpdateRequest) (*VirtualNetwork, *types.Response, error) {
	apiPath := path.Join(virtualNetworkBasePath, virtualNetworkID)
	virtualNetwork := new(VirtualNetworkGetResponse)

	resp, err := s.Client.DoRequest("PATCH", apiPath, updateRequest, virtualNetwork)
	if err != nil {
		return nil, resp, err
	}

	flatVirtualNetwork := NewFlatVirtualNetwork(virtualNetwork.Data)
	return &flatVirtualNetwork, resp, err
}

// Delete deletes a virtual network
func (s *VirtualNetworkServiceOp) Delete(virtualNetworkID string) (*types.Response, error) {
	apiPath := path.Join(virtualNetworkBasePath, virtualNetworkID)

	return s.Client.DoRequest("DELETE", apiPath, nil, nil)
}
