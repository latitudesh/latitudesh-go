package latitude

import (
	"encoding/json"
	"path"
	"strconv"
)

const planBasePath = "/plans"

// PlanService interface defines available plan methods
type PlanService interface {
	List(listOpt *ListOptions) ([]Plan, *Response, error)
	Get(string, *GetOptions) (*Plan, *Response, error)
}

// Plan represents a Latitude plan
type PlanRoot struct {
	Data PlanData `json:"data"`
	Meta meta     `json:"meta"`
}

type PlanListResponse struct {
	Data []PlanData `json:"data"`
	Meta meta       `json:"meta"`
}

type PlanData struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes PlanAttributes `json:"attributes"`
}

type PlanAttributes struct {
	Name          string             `json:"name"`
	Slug          string             `json:"slug"`
	Features      PlanFeatures       `json:"features"`
	Specs         PlanSpecs          `json:"specs"`
	Regions       PlanRegions        `json:"regions"`
	Availablility []PlanAvailability `json:"available_in"`
}

type PlanFeatures []string

type PlanSpecs struct {
	CPU    PlanCPU     `json:"cpu"`
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

type PlanRegions []PlanRegion

type PlanRegion struct {
	Name             string       `json:"name"`
	DeploysInstantly []string     `json:"deploys_instantly"`
	Locations        PlanLocation `json:"locations"`
	PlanPricing      Pricing      `json:"pricing"`
}

type PlanLocation struct {
	Available []string `json:"available"`
	InStock   []string `json:"in_stock"`
}

type Site struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	InStock    bool   `json:"in_stock"`
	StockLevel string `json:"stock_level"`
}

type Pricing struct {
	USD PricingUSD `json:"USD"`
	BRL PricingBRL `json:"BRL"`
}

type PricingUSD struct {
	Hour  float64 `json:"hour"`
	Month float64 `json:"month"`
	Year  float64 `json:"year"`
}

type PricingBRL struct {
	Hour  float64 `json:"hour"`
	Month float64 `json:"month"`
	Year  float64 `json:"year"`
}

type Plan struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Name     string       `json:"name"`
	Slug     string       `json:"slug"`
	Features PlanFeatures `json:"features"`
	Specs    PlanSpecs    `json:"specs"`
	InStock  []string     `json:"in_stock"`
}

// PlanServiceOp implements PlanService
type PlanServiceOp struct {
	client requestDoer
}

func (pd *PlanData) allStock() []string {
	stock := []string{}
	for _, region := range pd.Attributes.Regions {
		stock = append(stock, region.Locations.InStock...)
	}
	return stock
}

// Flatten latitude API data structures
func NewFlatPlan(pd PlanData) Plan {
	return Plan{
		pd.ID,
		pd.Type,
		pd.Attributes.Name,
		pd.Attributes.Slug,
		pd.Attributes.Features,
		pd.Attributes.Specs,
		pd.allStock(),
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
func (s *PlanServiceOp) List(opts *ListOptions) ([]Plan, *Response, error) {
	apiPathQuery := opts.WithQuery(planBasePath)

	plans := []Plan{}
	for {
		res := new(PlanListResponse)

		resp, err := s.client.DoRequest("GET", apiPathQuery, nil, res)
		if err != nil {
			return nil, resp, err
		}

		plans = append(plans, NewFlatPlanList(res.Data)...)

		if apiPathQuery = nextPage(res.Meta, opts); apiPathQuery != "" {
			continue
		}

		return plans, resp, err
	}
}

// Get returns a plan by id
func (s *PlanServiceOp) Get(planID string, opts *GetOptions) (*Plan, *Response, error) {
	endpointPath := path.Join(planBasePath, planID)
	apiPathQuery := opts.WithQuery(endpointPath)
	plan := new(PlanRoot)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, plan)
	if err != nil {
		return nil, resp, err
	}

	flatPlan := NewFlatPlan(plan.Data)
	return &flatPlan, resp, err
}
