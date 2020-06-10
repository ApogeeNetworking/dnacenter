package pnp

import (
	"bytes"
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
	Info           DeviceInfo `json:"deviceInfo"`
	WorkflowParams struct {
		CfgList          []DeviceConfig `json:"configList"`
		TopOfStackSerial string         `json:"topOfStackSerialNumber,omitempty"`
	} `json:"workflowParameters"`
	ID string `json:"id,omitempty"`
}

// FailedDevice ...
type FailedDevice struct {
	Index  int    `json:"index"`
	Serial string `json:"serialNum"`
	ID     string `json:"id"`
	Msg    string `json:"msg"`
}

// DeviceInfo is a PnP Device
type DeviceInfo struct {
	// Device Name in the DNA-C GUI
	// Cannot be more than 32 Chars Long
	Hostname string `json:"hostname"`
	// Fully Qualified Cisco Model
	ProductID string `json:"pid"`
	Serial    string `json:"serialNumber"`
	Stack     bool   `json:"stack"`
	// Unclaimed|Planned|Provisioned
	State        string    `json:"state,omitempty"`
	OnbState     string    `json:"onbState,omitempty"`
	ImageVersion string    `json:"imageVersion"`
	ProjectID    string    `json:"projectId,omitempty"`
	WorkflowID   string    `json:"workflowId,omitempty"`
	StackInfo    StackInfo `json:"stackInfo"`
}

// StackInfo ...
type StackInfo struct {
	// If TRUE then all Other Properties of this Struct will be omitted
	// Only IOS XE Versions SupportWorkFlows
	SupportsWorkflows bool          `json:"supportsStackWorkflows"`
	IsFullRing        bool          `json:"isFullRing,omitempty"`
	MemberList        []StackMember `json:"stackMemberList,omitempty"`
}

// StackMember ...
type StackMember struct {
	StackNumber int    `json:"stackNumber"`
	Serial      string `json:"serialNumber"`
	Role        string `json:"role"`
	Priority    int    `json:"priority"`
	MacAddr     string `json:"macAddress"`
	State       string `json:"state"`
	// The Fully Qualified Cisco Switch Model #
	ProductID string `json:"pid"`
}

// DCreds ...
type DCreds struct {
	User string `json:"username"`
	Pass string `json:"password"`
}

// GenResp for ResetDevice
type GenResp struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

// DevClaimResp ...
type DevClaimResp struct {
	Response string `json:"response"`
}

// BulkAddResp ...
type BulkAddResp struct {
	SuccessList []Device       `json:"successList"`
	FailureList []FailedDevice `json:"failureList"`
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
	j, _ = json.Marshal(res.Body)
	fmt.Println(string(j))
	return b, nil
}

// UpdateDevice ...
func (s *Service) UpdateDevice(device Device) {
	uri := fmt.Sprintf("%s/pnp-device/%s", s.baseURL, device.ID)
	d := Device{Info: DeviceInfo{Hostname: device.Info.Hostname}}
	j, _ := json.Marshal(d)
	body := strings.NewReader(string(j))
	res, err := s.http.MakeReq(uri, "PUT", body)
	if err != nil {
		// return b, fmt.Errorf("%v", err)
	}
	defer res.Body.Close()
	// j, _ = ioutil.ReadAll(res.Body)
	// fmt.Println(string(j))
}

// DeleteDevice ...
func (s *Service) DeleteDevice(id string) (string, error) {
	var device Device
	uri := fmt.Sprintf("%s/pnp-device/%s", s.baseURL, id)
	res, err := s.http.MakeReq(uri, "DELETE", nil)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	reader := ioutil.NopCloser(bytes.NewReader(data))
	defer reader.Close()
	if err = json.NewDecoder(reader).Decode(&device); err != nil {
		return "", fmt.Errorf("%v", err)
	}
	if device.ID == "" {
		reader = ioutil.NopCloser(bytes.NewReader(data))
		type resp struct {
			Response struct {
				Message string `json:"message"`
			} `json:"response"`
		}
		var response resp
		if err = json.NewDecoder(reader).Decode(&response); err != nil {
			return "", fmt.Errorf("%v", err)
		}
		return response.Response.Message, nil
	}
	return device.Info.State, nil
}

// GetDevicesBySerial ...
func (s *Service) GetDevicesBySerial(serials string) ([]Device, error) {
	uri := fmt.Sprintf("%s/pnp-device?serialNumber=%s", s.baseURL, serials)
	res, err := s.http.MakeReq(uri, "GET", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var devices []Device
	json.NewDecoder(res.Body).Decode(&devices)
	return devices, nil
}

// GetDevice ...
func (s *Service) GetDevice(id string) (Device, error) {
	uri := fmt.Sprintf("%s/pnp-device/%s", s.baseURL, id)
	res, err := s.http.MakeReq(uri, "GET", nil)
	if err != nil {
		return Device{}, fmt.Errorf("%v", err)
	}
	defer res.Body.Close()
	var device Device
	json.NewDecoder(res.Body).Decode(&device)
	return device, nil
}

func (s *Service) devRetrieval(devs []Device, offset int) []Device {
	var devices []Device
	uri := fmt.Sprintf("%s/pnp-device?offset=%v", s.baseURL, offset)
	res, _ := s.http.MakeReq(uri, "GET", nil)
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&devices)
	// Merged Device Slice(s)
	devs = append(devs, devices...)
	if len(devices) < offset {
		return devs
	}
	offset += 50
	return s.devRetrieval(devs, offset)
}

// GetDevices ...
func (s *Service) GetDevices() ([]Device, error) {
	return s.devRetrieval([]Device{}, 0), nil
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
	CableScheme string `json:"cablingScheme,omitempty"`
	ImageInfo   struct {
		ID   string `json:"imageId"`
		Skip bool   `json:"skip"`
	} `json:"imageInfo"`
	ConfigInfo DeviceConfig `json:"configInfo"`
}

// DeviceReset ...
type DeviceReset struct {
	DeviceList []DeviceList `json:"deviceResetList"`
	ProjectID  string       `json:"projectId,omitempty"`
	WorkflowID string       `json:"workflowId,omitempty"`
}

// DeviceList ...
type DeviceList struct {
	ConfigList       []DeviceConfig `json:"configList"`
	DeviceID         string         `json:"deviceId"`
	TopOfStackSerial string         `json:"topOfStackSerialNumber,omitempty"`
}

// DeviceConfig ...
type DeviceConfig struct {
	TemplateID string       `json:"configId"`
	Params     []TemplParam `json:"configParameters"`
}

// TemplParam ...
type TemplParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ClaimDeviceToSite ...
func (s *Service) ClaimDeviceToSite(sdc DeviceSiteClaim) DevClaimResp {
	var resp DevClaimResp
	uri := fmt.Sprintf("%s/pnp-device/site-claim", s.baseURL)
	j, _ := json.Marshal(sdc)
	body := strings.NewReader(string(j))
	res, err := s.http.MakeReq(uri, "POST", body)
	if err != nil {
		return resp
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&resp)
	return resp
}

// ResetDevice ...
func (s *Service) ResetDevice(dr DeviceReset) GenResp {
	var resp GenResp
	// https://dnac-ip/dna/intent/api/v1/onboarding + pnp-device/reset
	uri := fmt.Sprintf("%s/pnp-device/reset", s.baseURL)
	j, _ := json.Marshal(dr)
	body := strings.NewReader(string(j))
	res, err := s.http.MakeReq(uri, "POST", body)
	if err != nil {
		return resp
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&resp)
	return resp
}

// Settings PnP Settings
type Settings struct {
	TaskTimeouts   Timeouts `json:"taskTimeOuts"`
	TenantID       string   `json:"tenantId"`
	AAACredentails DCreds   `json:"aaaCredentials"`
	DefaultProfile Profile  `json:"defaultProfile"`
	AcceptEULA     bool     `json:"acceptEula"`
	ID             string   `json:"id"`
}

// Timeouts PnP Default Timeouts
type Timeouts struct {
	// All In Minutes
	ImageDownload int `json:"imageDownloadTimeOut"`
	Config        int `json:"configTimeOut"`
	General       int `json:"generalTimeOut"`
}

// Profile PnP Profile Settings
type Profile struct {
	Proxy   bool     `json:"proxy"`
	Cert    string   `json:"cert"`
	IPAddrs []string `json:"ipAddresses"`
	Port    int      `json:"port"`
}

// GetSettings ...
func (s *Service) GetSettings() (Settings, error) {
	var settings Settings
	uri := fmt.Sprintf("%s/pnp-settings", s.baseURL)
	res, err := s.http.MakeReq(uri, "GET", nil)
	if err != nil {
		return settings, err
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&settings)
	return settings, nil
}

// UpdateSettings ...
func (s *Service) UpdateSettings(settings Settings) (Settings, error) {
	oldSettings, _ := s.GetSettings()
	if settings.TaskTimeouts.Config == 0 {
		settings.TaskTimeouts.Config = oldSettings.TaskTimeouts.Config
	}
	if settings.TaskTimeouts.General == 0 {
		settings.TaskTimeouts.General = oldSettings.TaskTimeouts.General
	}
	if settings.TaskTimeouts.ImageDownload == 0 {
		settings.TaskTimeouts.ImageDownload = oldSettings.TaskTimeouts.ImageDownload
	}
	j, _ := json.Marshal(settings)
	body := strings.NewReader(string(j))
	res, err := s.http.MakeReq(s.baseURL+"/pnp-settings", "PUT", body)
	var newSettings Settings
	if err != nil {
		return newSettings, err
	}
	defer res.Body.Close()
	json.NewDecoder(res.Body).Decode(&newSettings)
	return newSettings, nil
}

// DeviceHistory ...
type DeviceHistory struct {
	TimeStamp int64           `json:"timestamp"`
	Details   string          `json:"details"`
	ErrorFlag bool            `json:"errorFlag"`
	TaskInfo  HistoryTaskInfo `json:"historyTaskInfo"`
}

// HistoryTaskInfo ...
type HistoryTaskInfo struct {
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	TimeTaken int        `json:"timeTaken"`
	List      []WorkItem `json:"workItemList"`
}

// WorkItem ...
type WorkItem struct {
	State        string `json:"state"`
	Command      string `json:"command"`
	AgentCommand string `json:"agentCommand"`
	CmdID        string `json:"cmdId"`
	RetryCount   int    `json:"retryCount"`
	StartTime    int64  `json:"startTime"`
	EndTime      int64  `json:"endTime"`
	TimeTaken    int    `json:"timeTaken"`
	Output       string `json:"outputStr"`
}

// GetDeviceHistory ...TaskTimeouts
func (s *Service) GetDeviceHistory(serial string) ([]DeviceHistory, error) {
	uri := fmt.Sprintf("%s/pnp-device/history?serialNumber=%s", s.baseURL, serial)
	type dhResp struct {
		Response []DeviceHistory `json:"response"`
	}
	var deviceHistory dhResp
	res, err := s.http.MakeReq(uri, "GET", nil)
	if err != nil {
		return deviceHistory.Response, err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&deviceHistory)
	if err != nil {
		return deviceHistory.Response, err
	}
	return deviceHistory.Response, nil
}
