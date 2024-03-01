package regions

import (
	"path"

	api "github.com/latitudesh/latitudesh-go/api_utils"
	internal "github.com/latitudesh/latitudesh-go/internal"
	types "github.com/latitudesh/latitudesh-go/types"
)

const regionBasePath = "/regions"

// RegionService interface defines available region methods
type RegionService interface {
	List(listOpt *api.ListOptions) ([]Region, *types.Response, error)
	Get(string, *api.GetOptions) (*Region, *types.Response, error)
}

// Plan represents a Latitude plan
type RegionRoot struct {
	Data RegionData    `json:"data"`
	Meta internal.Meta `json:"meta"`
}

type RegionListResponse struct {
	Data []RegionData  `json:"data"`
	Meta internal.Meta `json:"meta"`
}

type RegionData struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes RegionAttributes `json:"attributes"`
}

type RegionAttributes struct {
	Name     string        `json:"name"`
	Slug     string        `json:"slug"`
	Facility string        `json:"facility"`
	Country  RegionCountry `json:"country"`
}

type RegionCountry struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Region struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Facility    string `json:"facility"`
	CountryName string `json:"country_name"`
	CountrySlug string `json:"country_slug"`
}

// RegionServiceOp implements RegionService
type RegionServiceOp struct {
	Client internal.RequestDoer
}

// Flatten latitude API data structures
func NewFlatRegion(rd RegionData) Region {
	return Region{
		rd.ID,
		rd.Type,
		rd.Attributes.Name,
		rd.Attributes.Slug,
		rd.Attributes.Facility,
		rd.Attributes.Country.Name,
		rd.Attributes.Country.Slug,
	}
}

func NewFlatRegionList(rd []RegionData) []Region {
	var res []Region
	for _, region := range rd {
		res = append(res, NewFlatRegion(region))
	}
	return res
}

// List returns a list of regions
func (s *RegionServiceOp) List(opts *api.ListOptions) (regions []Region, resp *types.Response, err error) {
	apiPathQuery := opts.WithQuery(regionBasePath)

	for {
		res := new(RegionListResponse)

		resp, err = s.Client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		regions = append(regions, NewFlatRegionList(res.Data)...)

		if apiPathQuery = api.NextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

// Get returns a region by id
func (s *RegionServiceOp) Get(regionID string, opts *api.GetOptions) (*Region, *types.Response, error) {
	endpointPath := path.Join(regionBasePath, regionID)
	apiPathQuery := opts.WithQuery(endpointPath)
	region := new(RegionRoot)
	resp, err := s.Client.DoRequest("GET", apiPathQuery, nil, region)
	if err != nil {
		return nil, resp, err
	}

	flatRegion := NewFlatRegion(region.Data)
	return &flatRegion, resp, err
}
