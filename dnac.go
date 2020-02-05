package dnac

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const loginWarning string = "you must first login to perform this action"

// Client used for CDNAC Connection Handler
type Client struct {
	BaseURL  string
	Username string
	Password string
	IP       string
	// These are purposely not exported
	http      *http.Client
	authToken string
}

// NewClient creates a reference to the DNAC Client Struct
func NewClient(host, user, pass string, ignoreSSL bool) *Client {
	return &Client{
		BaseURL:  fmt.Sprintf("https://%s", host),
		Username: user,
		Password: pass,
		IP:       host,
		http: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: ignoreSSL,
				},
			},
			Timeout: 8 * time.Second,
		},
	}
}

func (c *Client) genReq(url, method string, r io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.BaseURL+url, r)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	if c.authToken != "" {
		req.Header.Set("X-AUTH-TOKEN", c.authToken)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}
	if strings.Contains(url, "/v1/auth") {
		req.SetBasicAuth(c.Username, c.Password)
	}

	return req, nil
}

// Login establishes a session with the DNAC API
func (c *Client) Login() error {
	// First we have to perform an Auth Login
	req, err := c.genReq("/api/system/v1/auth/login", "GET", nil)
	if err != nil {
		return err
	}
	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	// Grab the Cookie for the Next REQ
	// We Do Not Have to Do this Afterward
	authCookie := res.Cookies()[0]

	c.BaseURL = fmt.Sprintf("https://%s/dna", c.IP)

	req, err = c.genReq("/system/api/v1/auth/token", "POST", nil)
	if err != nil {
		return fmt.Errorf("failed to create a request: %v", err)
	}
	req.AddCookie(authCookie)

	res, err = c.http.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer res.Body.Close()
	// var tok AuthTok
	var tok struct {
		Token string `json:"Token"`
	}
	json.NewDecoder(res.Body).Decode(&tok)
	c.authToken = tok.Token
	return nil
}

// Device contains the properties of a Network Device
type Device struct {
	Hostname                  string      `json:"hostname"`
	RoleSource                string      `json:"roleSource"`
	InventoryStatusDetail     string      `json:"inventoryStatusDetail"`
	ApManagerInterfaceIP      string      `json:"apManagerInterfaceIp"`
	AssociatedWlcIP           string      `json:"associatedWlcIp"`
	BootDateTime              interface{} `json:"bootDateTime"`
	CollectionStatus          string      `json:"collectionStatus"`
	ErrorCode                 string      `json:"errorCode"`
	ErrorDescription          interface{} `json:"errorDescription"`
	Family                    string      `json:"family"`
	InterfaceCount            string      `json:"interfaceCount"`
	LastUpdated               string      `json:"lastUpdated"`
	LineCardCount             string      `json:"lineCardCount"`
	LineCardID                string      `json:"lineCardId"`
	LocationName              interface{} `json:"locationName"`
	ManagementIPAddress       string      `json:"managementIpAddress"`
	MemorySize                string      `json:"memorySize"`
	PlatformID                string      `json:"platformId"`
	ReachabilityFailureReason string      `json:"reachabilityFailureReason"`
	ReachabilityStatus        string      `json:"reachabilityStatus"`
	Series                    string      `json:"series"`
	SnmpContact               string      `json:"snmpContact"`
	SnmpLocation              string      `json:"snmpLocation"`
	TagCount                  string      `json:"tagCount"`
	TunnelUDPPort             string      `json:"tunnelUdpPort"`
	UptimeSeconds             int         `json:"uptimeSeconds"`
	WaasDeviceMode            interface{} `json:"waasDeviceMode"`
	CollectionInterval        string      `json:"collectionInterval"`
	SoftwareType              interface{} `json:"softwareType"`
	SoftwareVersion           string      `json:"softwareVersion"`
	LastUpdateTime            int64       `json:"lastUpdateTime"`
	UpTime                    string      `json:"upTime"`
	DeviceSupportLevel        string      `json:"deviceSupportLevel"`
	MacAddress                string      `json:"macAddress"`
	SerialNumber              string      `json:"serialNumber"`
	Type                      string      `json:"type"`
	Description               interface{} `json:"description"`
	Location                  interface{} `json:"location"`
	Role                      string      `json:"role"`
	InstanceTenantID          string      `json:"instanceTenantId"`
	InstanceUUID              string      `json:"instanceUuid"`
	ID                        string      `json:"id"`
}

// NetDeviceRes contains the Response of a DNAC Request
type NetDeviceRes struct {
	Response []Device `json:"response"`
	Version  string   `json:"version"`
}

// GetNetDevice retrieves Network Devices in DNAC
func (c *Client) GetNetDevice() (NetDeviceRes, error) {
	if c.authToken == "" {
		return NetDeviceRes{}, fmt.Errorf(loginWarning)
	}
	req, err := c.genReq("/intent/api/v1/network-device", "GET", nil)
	if err != nil {
		return NetDeviceRes{}, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return NetDeviceRes{}, fmt.Errorf("request failed: %v", err)
	}
	defer res.Body.Close()
	var netDevice NetDeviceRes
	json.NewDecoder(res.Body).Decode(&netDevice)
	return netDevice, nil
}
