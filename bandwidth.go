package latitude

type BandwidthService interface {
	TrafficQuota(opts *ListOptions) (*TrafficQuota, *Response, error)
	TrafficConsumption(opts *ListOptions) (*TrafficConsumption, *Response, error)
}

// TrafficConsumption represents consumed traffic in regions
type TrafficConsumption struct {
	FromDate                        int64           `json:"from_date"`
	ToDate                          int64           `json:"to_date"`
	TotalInboundGB                  int64           `json:"total_inbound_gb"`
	TotalOutboundGB                 int64           `json:"total_outbound_gb"`
	TotalInbound95thPercentileMbps  float64         `json:"total_inbound_95th_percentile_mbps"`
	TotalOutbound95thPercentileMbps float64         `json:"total_outbound_95th_percentile_mbps"`
	Regions                         []TrafficRegion `json:"regions"`
}

type TrafficQuota struct {
	QuotaPerProject []QuotaPerProject `json:"quota_per_project"`
}

type QuotaPerProject struct {
	ProjectID      int64            `json:"project_id"`
	ProjectSlug    string           `json:"project_slug"`
	BillingMethod  string           `json:"billing_method"`
	QuotaPerRegion []QuotaPerRegion `json:"quota_per_region"`
}

type QuotaPerRegion struct {
	RegionId    int64       `json:"region_id"`
	RegionSlug  string      `json:"region_slug"`
	QuotaInTb   QuotaDetail `json:"quota_in_tb"`
	QuotaInMbps QuotaDetail `json:"quota_in_mbps"`
}
type QuotaDetail struct {
	Granted    int64 `json:"granted"`
	Additional int64 `json:"additional"`
	Total      int64 `json:"total"`
}

type TrafficRegion struct {
	RegionSlug string              `json:"region_slug"`
	Data       []TrafficRegionData `json:"data"`
}

type TrafficRegionData struct {
	Date                 string  `json:"date"`
	InboundGB            int64   `json:"inbound_gb"`
	OutboundGB           int64   `json:"outbound_gb"`
	AvgOutboundSpeedMbps float64 `json:"avg_outbound_speed_mbps"`
	AvgInboundSpeedMbps  float64 `json:"avg_inbound_speed_mbps"`
	OutboundSpeedMbps    float64 `json:"outbound_speed_mbps"`
	InboundSpeedMbps     float64 `json:"inbound_speed_mbps"`
}

// BandwidthServiceOp implements BandwidthService
type BandwidthServiceOp struct {
	client requestDoer
}

type TrafficConsumptionResponse struct {
	Data TrafficConsumptionData `json:"data"`
	Meta meta                   `json:"meta"`
}

type TrafficQuotaResponse struct {
	Data TrafficQuotaData `json:"data"`
	Meta meta             `json:"meta"`
}

type TrafficConsumptionData struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Attributes TrafficConsumption `json:"attributes"`
}

type TrafficQuotaData struct {
	ID         string       `json:"id"`
	Type       string       `json:"type"`
	Attributes TrafficQuota `json:"attributes"`
}

// Flatten latitude API data structures
func NewFlatTrafficConsumption(t TrafficConsumptionData) TrafficConsumption {
	return TrafficConsumption{
		t.Attributes.FromDate,
		t.Attributes.ToDate,
		t.Attributes.TotalInboundGB,
		t.Attributes.TotalOutboundGB,
		t.Attributes.TotalInbound95thPercentileMbps,
		t.Attributes.TotalOutbound95thPercentileMbps,
		t.Attributes.Regions,
	}
}

func NewFlatTrafficQuota(t TrafficQuotaData) TrafficQuota {
	return TrafficQuota{
		t.Attributes.QuotaPerProject,
	}
}

// TrafficConsumption returns consumed traffic
func (u *BandwidthServiceOp) TrafficConsumption(opts *ListOptions) (*TrafficConsumption, *Response, error) {
	trafficConsumptionResponse := new(TrafficConsumptionResponse)
	apiPathQuery := opts.WithQuery("/traffic")

	resp, err := u.client.DoRequest("GET", apiPathQuery, nil, trafficConsumptionResponse)
	if err != nil {
		return nil, resp, err
	}

	flatFlatTrafficConsumption := NewFlatTrafficConsumption(trafficConsumptionResponse.Data)
	return &flatFlatTrafficConsumption, resp, err
}

// TrafficQuota returns purchased quota
func (u *BandwidthServiceOp) TrafficQuota(opts *ListOptions) (*TrafficQuota, *Response, error) {
	trafficQuotaResponse := new(TrafficQuotaResponse)
	apiPathQuery := opts.WithQuery("/traffic/quota")

	resp, err := u.client.DoRequest("GET", apiPathQuery, nil, trafficQuotaResponse)
	if err != nil {
		return nil, resp, err
	}

	flatFlatTrafficQuota := NewFlatTrafficQuota(trafficQuotaResponse.Data)
	return &flatFlatTrafficQuota, resp, err
}
