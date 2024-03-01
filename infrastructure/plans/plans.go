package plans

import (
	"encoding/json"
	"path"
	"strconv"

	api "github.com/latitudesh/latitudesh-go/api_utils"
	internal "github.com/latitudesh/latitudesh-go/internal"
	types "github.com/latitudesh/latitudesh-go/types"
)

const planBasePath = "/plans"

// PlanService interface defines available plan methods
type PlanService interface {
	List(listOpt *api.ListOptions) ([]Plan, *types.Response, error)
	Get(string, *api.GetOptions) (*Plan, *types.Response, error)
}

// Plan represents a Latitude plan
type PlanRoot struct {
	Data PlanData      `json:"data"`
	Meta internal.Meta `json:"meta"`
}

type PlanListResponse struct {
	Data []PlanData    `json:"data"`
	Meta internal.Meta `json:"meta"`
}

type PlanData struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes PlanAttributes `json:"attributes"`
}

type PlanAttributes struct {
	Name          string             `json:"name"`
	Slug          string             `json:"slug"`
	Line          string             `json:"line"`
	Features      PlanFeatures       `json:"features"`
	Specs         PlanSpecs          `json:"specs"`
	Availablility []PlanAvailability `json:"available_in"`
}

type PlanFeatures []string

type PlanSpecs struct {
	CPUs   []PlanCPU   `json:"cpus"`
	Memory PlanMemory  `json:"memory"`
	Drives []PlanDrive `json:"drives"`
	NICs   []PlanNIC   `json:"nics"`
}

type PlanCPU struct {
	Type  string  `json:"type"`
	Clock float64 `json:"clock"`
	Cores int     `json:"cores"`
	Count int     `json:"count"`
}

type PlanMemory struct {
	Total string `json:"total"`
}

func (p *PlanMemory) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var total string
	err = json.Unmarshal(*objMap["total"], &total)
	if err != nil {
		// total is an int
		var totalInt int
		err = json.Unmarshal(*objMap["total"], &totalInt)
		if err != nil {
			return err
		}
		ts := strconv.Itoa(totalInt)
		p.Total = ts
	} else {
		p.Total = total
	}

	return nil
}

type PlanDrive struct {
	Count int    `json:"count"`
	Size  string `json:"size"`
	Type  string `json:"type"`
}

type PlanNIC struct {
	Count int    `json:"count"`
	Type  string `json:"type"`
}

type PlanAvailability struct {
	Region  PlanRegion `json:"region"`
	Sites   []Site     `json:"sites"`
	Pricing Pricing    `json:"pricing"`
}

type PlanRegion struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Site struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	InStock    bool   `json:"in_stock"`
	StockLevel string `json:"stock_level"`
}

type Pricing struct {
	USD PricingUSD `json:"USD"`
	BRL PricingBRL `json:"name"`
}

type PricingUSD struct {
	Month float64 `json:"month"`
}

type PricingBRL struct {
	Month float64 `json:"month"`
}

type Plan struct {
	ID           string             `json:"id"`
	Type         string             `json:"type"`
	Name         string             `json:"name"`
	Slug         string             `json:"slug"`
	Line         string             `json:"line"`
	Features     PlanFeatures       `json:"features"`
	Specs        PlanSpecs          `json:"specs"`
	Availibility []PlanAvailability `json:"availibility"`
}

// PlanServiceOp implements PlanService
type PlanServiceOp struct {
	Client internal.RequestDoer
}

// Flatten latitude API data structures
func NewFlatPlan(pd PlanData) Plan {
	return Plan{
		pd.ID,
		pd.Type,
		pd.Attributes.Name,
		pd.Attributes.Slug,
		pd.Attributes.Line,
		pd.Attributes.Features,
		pd.Attributes.Specs,
		pd.Attributes.Availablility,
	}
}

func NewFlatPlanList(pd []PlanData) []Plan {
	var res []Plan
	for _, plan := range pd {
		res = append(res, NewFlatPlan(plan))
	}
	return res
}

// List returns a list of plans
func (s *PlanServiceOp) List(opts *api.ListOptions) (plans []Plan, resp *types.Response, err error) {
	apiPathQuery := opts.WithQuery(planBasePath)

	for {
		res := new(PlanListResponse)

		resp, err = s.Client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		plans = append(plans, NewFlatPlanList(res.Data)...)

		if apiPathQuery = api.NextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

// Get returns a plan by id
func (s *PlanServiceOp) Get(planID string, opts *api.GetOptions) (*Plan, *types.Response, error) {
	endpointPath := path.Join(planBasePath, planID)
	apiPathQuery := opts.WithQuery(endpointPath)
	plan := new(PlanRoot)
	resp, err := s.Client.DoRequest("GET", apiPathQuery, nil, plan)
	if err != nil {
		return nil, resp, err
	}

	flatPlan := NewFlatPlan(plan.Data)
	return &flatPlan, resp, err
}
