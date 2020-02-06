package dnac

import (
	"encoding/json"
	"fmt"
)

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

// DeviceRes contains the Response of a DNAC Request
type DeviceRes struct {
	Response []Device `json:"response"`
	Version  string   `json:"version"`
}

// GetNetDevice retrieves Network Devices in DNAC
func (c *Client) GetNetDevice() (DeviceRes, error) {
	if c.authToken == "" {
		return DeviceRes{}, fmt.Errorf(loginWarning)
	}
	res, err := c.MakeReq("/intent/api/v1/network-device", "GET", nil)
	if err != nil {
		return DeviceRes{}, nil
	}
	defer res.Body.Close()
	var netDevice DeviceRes
	json.NewDecoder(res.Body).Decode(&netDevice)
	return netDevice, nil
}

// DeviceVLAN contains the VLAN Props of a Device
type DeviceVLAN struct {
	VlanNumber     int    `json:"vlanNumber"`
	NumberOfIPs    int    `json:"numberOfIPs,omitempty"`
	IPAddress      string `json:"ipAddress,omitempty"`
	Prefix         string `json:"prefix,omitempty"`
	InterfaceName  string `json:"interfaceName"`
	NetworkAddress string `json:"networkAddress,omitempty"`
}

// DeviceVlanRes contains the resp of the Request
type DeviceVlanRes struct {
	Response []DeviceVLAN `json:"response"`
	Version  string       `json:"version"`
}

// GetDeviceVLANs retrieve VLANs associated with a Device
func (c *Client) GetDeviceVLANs(id string) (DeviceVlanRes, error) {
	if c.authToken == "" {
		return DeviceVlanRes{}, fmt.Errorf(loginWarning)
	}
	endpoint := fmt.Sprintf("/intent/api/v1/network-device/%s/vlan", id)
	res, err := c.MakeReq(endpoint, "GET", nil)
	if err != nil {
		return DeviceVlanRes{}, err
	}
	defer res.Body.Close()
	var deviceVlans DeviceVlanRes
	json.NewDecoder(res.Body).Decode(&deviceVlans)
	return deviceVlans, nil
}
