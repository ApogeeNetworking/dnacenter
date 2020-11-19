package devices

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ApogeeNetworking/dnacenter/requests"
)

// Service ...
type Service struct {
	baseURL string
	http    *requests.Req
}

// New creates an instance of a DNA-C Device Service
func New(uri string, c *requests.Req) *Service {
	uri = strings.ReplaceAll(uri, "/dna", "")
	return &Service{
		baseURL: uri + "/api/v1/network-device",
		http:    c,
	}
}

// Get retrieves Network Devices in DNAC Inventory
func (s *Service) Get(param ReqParams) ([]Device, error) {
	var uri string
	switch {
	case param.Hostname != "":
		uri = fmt.Sprintf("%s?hostname=%s", s.baseURL, param.Hostname)
	case param.IPAddr != "":
		uri = fmt.Sprintf("%s?managementIpAddress=%s", s.baseURL, param.IPAddr)
	case param.Serial != "":
		uri = fmt.Sprintf("%s?serialNumber=%s", s.baseURL, param.Serial)
	default:
		uri = s.baseURL
	}
	res, err := s.http.MakeReq(uri, "GET", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resp := struct {
		Response []Device `json:"response"`
	}{}
	json.NewDecoder(res.Body).Decode(&resp)
	return resp.Response, nil
}

// GetByID ...
func (s *Service) GetByID(id string) (Device, error) {
	res, err := s.http.MakeReq(
		fmt.Sprintf("%s/%s", s.baseURL, id),
		"GET",
		nil,
	)
	if err != nil {
		return Device{}, err
	}
	defer res.Body.Close()
	resp := struct {
		Response Device `json:"response"`
	}{}
	json.NewDecoder(res.Body).Decode(&resp)
	return resp.Response, nil
}

// Delete ...
func (s *Service) Delete(id string, cleanCfg bool) (int, error) {
	res, err := s.http.MakeReq(
		fmt.Sprintf("%s/%s?cleanConfig=%v", s.baseURL, id, cleanCfg),
		"DELETE",
		nil,
	)
	if err != nil {
		return 0, err
	}
	return res.StatusCode, nil
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
