package devices

// Device contains the properties of a Network Device
type Device struct {
	ID                        string      `json:"id"`
	Hostname                  string      `json:"hostname"`
	IPAddr                    string      `json:"managementIpAddress"`
	MacAddr                   string      `json:"macAddress"`
	Serial                    string      `json:"serialNumber"`
	Type                      string      `json:"type"`
	Description               interface{} `json:"description"`
	LocationName              interface{} `json:"locationName"`
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
	Location                  interface{} `json:"location"`
	Role                      string      `json:"role"`
	InstanceTenantID          string      `json:"instanceTenantId"`
	InstanceUUID              string      `json:"instanceUuid"`
}

// VLAN contains the VLAN Props of a Device
type VLAN struct {
	VlanNumber     int    `json:"vlanNumber"`
	NumberOfIPs    int    `json:"numberOfIPs,omitempty"`
	IPAddress      string `json:"ipAddress,omitempty"`
	Prefix         string `json:"prefix,omitempty"`
	InterfaceName  string `json:"interfaceName"`
	NetworkAddress string `json:"networkAddress,omitempty"`
}

// ReqParams ...
type ReqParams struct {
	Hostname string
	IPAddr   string
	Serial   string
}
