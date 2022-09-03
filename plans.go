package latitude

import (
	"encoding/json"
	"path"
)

const planBasePath = "/plans"

// PlanService interface defines available plan methods
type PlanService interface {
	List(listOpt *ListOptions) ([]PlanData, *Response, error)
	Get(string, *GetOptions) (*Plan, *Response, error)
}

// Plan represents a Latitude plan
type Plan struct {
	Data PlanData `json:"data"`
	Meta meta     `json:"meta"`
}

// Plan represents a Latitude plan
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
	Name      string             `json:"name"`
	Slug      string             `json:"slug"`
	Line      string             `json:"line"`
	Features  PlanFeatures       `json:"features"`
	Specs     PlanSpecs          `json:"specs"`
	Available []PlanAvailability `json:"available_in"`
}

type PlanFeatures struct {
	SSH      bool `json:"ssh"`
	RAID     bool `json:"raid"`
	UserData bool `json:"user_data"`
}

type PlanSpecs struct {
	Name   string      `json:"name"`
	Slug   string      `json:"slug"`
	Line   string      `json:"line"`
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
	// Sometime total is return as a string and sometimes as an int
	Total json.Number `json:"total"`
}

type PlanDrive struct {
	Count int    `json:"total"`
	Size  string `json:"size"`
	Type  string `json:"type"`
}

type PlanNIC struct {
	Count int    `json:"total"`
	Type  string `json:"type"`
}

type PlanAvailability struct {
	Region  Region  `json:"region"`
	Sites   []Site  `json:"site"`
	Pricing Pricing `json:"pricing"`
}

type Region struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Site struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	InStock    bool   `json:"in_stock"`
	StockLevel string `json:"stock+_level"`
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

// PlanServiceOp implements PlanService
type PlanServiceOp struct {
	client requestDoer
}

// List returns a list of plans
func (s *PlanServiceOp) List(opts *ListOptions) (plans []PlanData, resp *Response, err error) {
	apiPathQuery := opts.WithQuery(planBasePath)

	for {
		subset := new(PlanListResponse)

		resp, err = s.client.DoRequest("GET", apiPathQuery, nil, subset)
		if err != nil {
			return nil, resp, err
		}

		plans = append(plans, subset.Data...)

		if apiPathQuery = nextPage(subset.Meta, opts); apiPathQuery != "" {
			continue
		}

		return
	}
}

// Get returns a plan by id
func (s *PlanServiceOp) Get(planID string, opts *GetOptions) (*Plan, *Response, error) {
	endpointPath := path.Join(planBasePath, planID)
	apiPathQuery := opts.WithQuery(endpointPath)
	plan := new(Plan)
	resp, err := s.client.DoRequest("GET", apiPathQuery, nil, plan)
	if err != nil {
		return nil, resp, err
	}
	return plan, resp, err
}
