package pnp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/drkchiloll/dnacenter/requests"
)

// Service ...
type Service struct {
	baseURL string
	http    *requests.Req
}

// New creates an instance of a DNA-C Plug N Play Service
func New(uri string, r *requests.Req) *Service {
	return &Service{baseURL: uri + "/intent/api/v1/onboarding", http: r}
}

// Device ...
type Device struct {
	Info DeviceInfo `json:"deviceInfo"`
	ID   string     `json:"id,omitempty"`
}

// DeviceInfo is a PnP Device
type DeviceInfo struct {
	Hostname  string `json:"hostname"`
	ProductID string `json:"pid"`
	Serial    string `json:"serialNumber"`
	Stack     bool   `json:"stack"`
	State     string `json:"state,omitempty"`
	OnbState  string `json:"onbState,omitempty"`
	// AAACreds       DCreds        `json:"aaaCredentials"`
	// SudiRequired   bool          `json:"sudiRequired"`
	// UserSudiSerial []string `json:"userSudiSerialNos"`
}

// DCreds ...
type DCreds struct {
	User string `json:"username"`
	Pass string `json:"password"`
}

// BulkAddResp ...
type BulkAddResp struct {
	SuccessList []Device `json:"successList"`
	FailureList []Device `json:"failureList"`
}

// BulkAddDevices to the PnP Database
func (s *Service) BulkAddDevices(devices []Device) (BulkAddResp, error) {
	var b BulkAddResp
	// Marshal devices into JSON Object
	j, _ := json.Marshal(devices)
	body := strings.NewReader(string(j))
	res, err := s.http.MakeReq(
		fmt.Sprintf("%s/pnp-device/import", s.baseURL),
		"POST",
		body,
	)
	if err != nil {
		return b, fmt.Errorf("%v", err)
	}
	err = json.NewDecoder(res.Body).Decode(&b)
	if err != nil {
		return b, fmt.Errorf("%v", err)
	}
	return b, nil
}

// GetDevice ...
func (s *Service) GetDevice(id string) (Device, error) {
	uri := fmt.Sprintf("%s/pnp-device/%s", s.baseURL, id)
	res, err := s.http.MakeReq(uri, "GET", nil)
	if err != nil {
	}
	defer res.Body.Close()
	var device Device
	json.NewDecoder(res.Body).Decode(&device)
	return device, nil
}

type cableScheme struct {
	OneA string
	OneB string
}

type switchType struct {
	Default       string
	StackSwitch   string
	AP            string
	CatWLC        string
	Sensor        string
	MobileExpress string
	SiteProvReq   string
}

var (
	// XeStackCabling Cisco IOS XE Stack Cabling Scheme
	XeStackCabling = cableScheme{
		// 1A Stack Cabling Means Port 1 of the Master Switch (top)
		// Is Connected to Port 2 of the Bottom Switch (last switch)
		OneA: "1A",
		// 1B Stack Cabling is the Inverse of 1A
		// If the Cabling Scheme is Reversed then Cisco DNA-C will Renumber
		// The Switches which could result in the Master Switch Being the Bottom Switch
		OneB: "1B",
	}
	// DeviceClaimType ...
	DeviceClaimType = switchType{
		Default:       "Default",
		StackSwitch:   "StackSwitch",
		AP:            "AccessPoint",
		CatWLC:        "CatalystWLC",
		Sensor:        "Sensor",
		MobileExpress: "MobilityExpress",
		SiteProvReq:   "SiteProvisionRequest",
	}
)

// DeviceSiteClaim ...
type DeviceSiteClaim struct {
	// Site defined in Network Profile (siteprofile endpoint)
	SiteID string `json:"siteId"`
	// PnP Device ID
	DeviceID string `json:"deviceId"`
	// PnP Device Hostname
	Hostname string `json:"hostname"`
	// Default|SwitchStack
	Type string `json:"type"`
	// Needed for IOS-XE Device (3850|9200L)
	TopOfStackSerial string `json:"topOfStackSerialNumber,omitempty"`
	// Needed for IOS-XE Device (3850|9200L)
	// 1A|1B
	CableScheme string `json:"cablingScheme"`
	ImageInfo   struct {
		ID   string `json:"imageId"`
		Skip bool   `json:"skip"`
	} `json:"imageInfo"`
	ConfigInfo struct {
		// TemplateID
		ID string `json:"configId"`
		// Template Variables with Values Specified
		Params []TemplParam `json:"configParameters"`
	} `json:"configInfo"`
}

// TemplParam ...
type TemplParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ClaimDeviceToSite ...
func (s *Service) ClaimDeviceToSite(sdc DeviceSiteClaim) {
	uri := fmt.Sprintf("%s/pnp-device/site-claim", s.baseURL)
	j, _ := json.Marshal(sdc)
	body := strings.NewReader(string(j))
	res, err := s.http.MakeReq(uri, "POST", body)
	if err != nil {
	}
	defer res.Body.Close()
	j, _ = ioutil.ReadAll(res.Body)
	fmt.Println(string(j))
}
