package latitude

const operatingSystemBasePath = "/plans/operating_systems"

// RegionService interface defines available region methods
type OperatinSystemService interface {
	List(listOpt *ListOptions) ([]Region, *Response, error)
}

type OperatingSystemListResponse struct {
	Data []OperatingSystemData `json:"data"`
	Meta meta                  `json:"meta"`
}

type OperatingSystemData struct {
	ID         string                    `json:"id"`
	Type       string                    `json:"type"`
	Attributes OperatingSystemAttributes `json:"attributes"`
}

type OperatingSystemAttributes struct {
	Name     string                  `json:"name"`
	Distro   string                  `json:"distro"`
	Slug     string                  `json:"slug"`
	Version  string                  `json:"version"`
	User     string                  `json:"user"`
	Features OperatingSystemFeatures `json:"features"`
}

type OperatingSystemFeatures struct {
	Raid     bool `json:"raid"`
	Rescue   bool `json:"rescue"`
	SshKeys  bool `json:"ssh_keys"`
	UserData bool `json:"user_data"`
}

type OperatingSystem struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Distro   string `json:"distro"`
	Slug     string `json:"slug"`
	Version  string `json:"version"`
	User     string `json:"user"`
	Raid     bool   `json:"raid"`
	Rescue   bool   `json:"rescue"`
	SshKeys  bool   `json:"ssh_keys"`
	UserData bool   `json:"user_data"`
}

type OperatingSystemServiceOp struct {
	client requestDoer
}

func NewFlatOperatingSystem(osd OperatingSystemData) OperatingSystem {
	return OperatingSystem{
		osd.ID,
		osd.Type,
		osd.Attributes.Name,
		osd.Attributes.Distro,
		osd.Attributes.Slug,
		osd.Attributes.Version,
		osd.Attributes.User,
		osd.Attributes.Features.Raid,
		osd.Attributes.Features.Rescue,
		osd.Attributes.Features.SshKeys,
		osd.Attributes.Features.UserData,
	}
}

func NewFlatOperatingSystemList(osd []OperatingSystemData) []OperatingSystem {
	var res []OperatingSystem
	for _, os := range osd {
		res = append(res, NewFlatOperatingSystem(os))
	}
	return res
}

// List returns a list of regions
func (os *OperatingSystemServiceOp) List(opts *ListOptions) (operatingSystems []OperatingSystem, resp *Response, err error) {
	apiPathQuery := opts.WithQuery(operatingSystemBasePath)

	for {
		res := new(OperatingSystemListResponse)

		resp, err = os.client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		operatingSystems = append(operatingSystems, NewFlatOperatingSystemList(res.Data)...)

		if apiPathQuery = nextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}
		return
	}
}
