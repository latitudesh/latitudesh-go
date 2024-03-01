package latitude

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	internal "github.com/latitudesh/latitudesh-go/internal"
	types "github.com/latitudesh/latitudesh-go/types"

	bandwidth "github.com/latitudesh/latitudesh-go/bandwidth"
	members "github.com/latitudesh/latitudesh-go/members"
	opsys "github.com/latitudesh/latitudesh-go/operating_systems"
	plans "github.com/latitudesh/latitudesh-go/plans"
	projects "github.com/latitudesh/latitudesh-go/projects"
	regions "github.com/latitudesh/latitudesh-go/regions"
	servers "github.com/latitudesh/latitudesh-go/servers"
	sshkeys "github.com/latitudesh/latitudesh-go/ssh_keys"
	teams "github.com/latitudesh/latitudesh-go/teams"
	userdata "github.com/latitudesh/latitudesh-go/user_data"
	vnet "github.com/latitudesh/latitudesh-go/virtual_networks"
	vlanassign "github.com/latitudesh/latitudesh-go/virtual_networks_assignments"
)

const (
	authTokenEnvVar      = "LATITUDE_AUTH_TOKEN"
	apiVersion           = "2023-06-01"
	baseURL              = "https://api.latitude.sh"
	debugEnvVar          = "LATITUDE_DEBUG"
	userAgentForSDK      = "Latitude-Go-SDK"
	userAgentForProvider = "Latitude-Terraform-Provider"
)

var currentVersion = "0.3.1"

// Client is the base API Client
type Client struct {
	client        *http.Client
	debug         bool
	BaseURL       *url.URL
	UserAgent     string
	ConsumerToken string
	APIKey        string

	Projects         projects.ProjectService
	Servers          servers.ServerService
	UserData         userdata.UserDataService
	SSHKeys          sshkeys.SSHKeyService
	Plans            plans.PlanService
	OperatingSystems opsys.OperatingSystemService
	VirtualNetworks  vnet.VirtualNetworkService
	VlanAssignments  vlanassign.VlanAssignmentService
	Regions          regions.RegionService
	Teams            teams.TeamService
	Bandwidth        bandwidth.BandwidthService
	Members          members.MemberService
}

// NewRequest inits a new http request with the proper headers
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	// relative path to append to the endpoint url, no leading slash please
	if path[0] == '/' {
		path = path[1:]
	}
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	// json encode the request body, if any
	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Close = true

	req.Header.Add("Authorization", c.APIKey)

	req.Header.Add("API-Version", apiVersion)

	// set User-Agent value for SDK or terraform-provider
	userAgent := c.UserAgent
	if !strings.Contains(userAgent, userAgentForProvider) {
		userAgent = fmt.Sprintf("%s/%s", userAgentForSDK, currentVersion)
	}
	req.Header.Add("User-Agent", userAgent)

	if req.Method != "GET" {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

// Do executes the http request
func (c *Client) Do(req *http.Request, v interface{}) (*types.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := types.Response{Response: resp}

	if c.debug {
		internal.DumpResponse(response.Response)
	}
	internal.DumpDeprecation(response.Response)

	err = internal.CheckResponse(resp)
	// if the response is an error, return the ErrorResponse
	if err != nil {
		return &response, err
	}

	if v != nil {
		// if v implements the io.Writer interface, return the raw response
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return &response, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return &response, err
			}
		}
	}

	return &response, err
}

// DoRequest is a convenience method, it calls NewRequest followed by Do
// v is the interface to unmarshal the response JSON into
func (c *Client) DoRequest(method, path string, body, v interface{}) (*types.Response, error) {
	req, err := c.NewRequest(method, path, body)
	if c.debug {
		internal.DumpRequest(req)
	}
	if err != nil {
		return nil, err
	}
	return c.Do(req, v)
}

// DoRequestWithHeader same as DoRequest
func (c *Client) DoRequestWithHeader(method string, headers map[string]string, path string, body, v interface{}) (*types.Response, error) {
	req, err := c.NewRequest(method, path, body)
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	if c.debug {
		internal.DumpRequest(req)
	}
	if err != nil {
		return nil, err
	}
	return c.Do(req, v)
}

// NewClient initializes and returns a Client
func NewClient() (*Client, error) {
	apiToken := os.Getenv(authTokenEnvVar)
	if apiToken == "" {
		return nil, fmt.Errorf("you must export %s", authTokenEnvVar)
	}
	c := NewClientWithAuth("latitude lib", apiToken, nil)
	return c, nil

}

// NewClientWithAuth initializes and returns a Client
func NewClientWithAuth(consumerToken string, apiKey string, httpClient *http.Client) *Client {
	client, _ := NewClientWithBaseURL(apiKey, httpClient, baseURL)
	return client
}

// NewClientWithBaseURL returns a Client pointing to nonstandard API URL, e.g.
// for mocking the remote API
func NewClientWithBaseURL(apiKey string, httpClient *http.Client, apiBaseURL string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	u, err := url.Parse(apiBaseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{client: httpClient, BaseURL: u, APIKey: apiKey}
	c.Projects = &projects.ProjectServiceOp{Client: c}
	c.Servers = &servers.ServerServiceOp{Client: c}
	c.SSHKeys = &sshkeys.SSHKeyServiceOp{Client: c}
	c.UserData = &userdata.UserDataServiceOp{Client: c}
	c.Teams = &teams.TeamServiceOp{Client: c}
	c.Bandwidth = &bandwidth.BandwidthServiceOp{Client: c}
	c.Plans = &plans.PlanServiceOp{Client: c}
	c.OperatingSystems = &opsys.OperatingSystemServiceOp{Client: c}
	c.Regions = &regions.RegionServiceOp{Client: c}
	c.VirtualNetworks = &vnet.VirtualNetworkServiceOp{Client: c}
	c.VlanAssignments = &vlanassign.VlanAssignmentServiceOp{Client: c}
	c.Members = &members.MemberServiceOp{Client: c}
	c.debug = os.Getenv(debugEnvVar) != ""

	return c, nil
}
