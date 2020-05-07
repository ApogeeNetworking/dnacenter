package devices

import (
	"encoding/json"
	"fmt"

	"github.com/drkchiloll/dnacenter/requests"
)

// Service ...
type Service struct {
	baseURL string
	http    *requests.Req
}

// New creates an instance of a DNA-C Device Service
func New(uri string, c *requests.Req) *Service {
	return &Service{
		baseURL: uri + "/intent/api/v1/network-device",
		http:    c,
	}
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

// Resp contains the Response of a DNAC Request
type Resp struct {
	Response []Device `json:"response"`
	Version  string   `json:"version"`
}

// Get retrieves Network Devices in DNAC Inventory
func (s *Service) Get() (Resp, error) {
	res, err := s.http.MakeReq(s.baseURL, "GET", nil)
	if err != nil {
		return Resp{}, nil
	}
	defer res.Body.Close()
	var netDevice Resp
	json.NewDecoder(res.Body).Decode(&netDevice)
	return netDevice, nil
}

// Delete ...
func (s *Service) Delete(id string, cleanCfg bool) {}

// VLAN contains the VLAN Props of a Device
type VLAN struct {
	VlanNumber     int    `json:"vlanNumber"`
	NumberOfIPs    int    `json:"numberOfIPs,omitempty"`
	IPAddress      string `json:"ipAddress,omitempty"`
	Prefix         string `json:"prefix,omitempty"`
	InterfaceName  string `json:"interfaceName"`
	NetworkAddress string `json:"networkAddress,omitempty"`
}

// VlanRes contains the resp of the Request
type VlanRes struct {
	Response []VLAN `json:"response"`
	Version  string `json:"version"`
}

// GetDeviceVLANs retrieve VLANs associated with a Device
func (s *Service) GetDeviceVLANs(id string) (VlanRes, error) {
	endpoint := fmt.Sprintf("%s/%s/vlan", s.baseURL, id)
	res, err := s.http.MakeReq(endpoint, "GET", nil)
	if err != nil {
		return VlanRes{}, err
	}
	defer res.Body.Close()
	var deviceVlans VlanRes
	json.NewDecoder(res.Body).Decode(&deviceVlans)
	return deviceVlans, nil
}
