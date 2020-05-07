package siteprofile

import (
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/drkchiloll/dnacenter/requests"
)

// Service ...package
type Service struct {
	baseURL string
	http    *requests.Req
}

// New creates an instance of a DNA-C NETWORK-PROFILE Service
func New(uri string, r *requests.Req) *Service {
	parsedURL, _ := url.Parse(uri)
	uriHost := parsedURL.Host
	ep := fmt.Sprintf("https://%s/api/v1/siteprofile", uriHost)
	return &Service{baseURL: ep, http: r}
}

// Get ...
func (s *Service) Get() {
	res, err := s.http.MakeReq(s.baseURL, "GET", nil)
	if err != nil {
	}
	defer res.Body.Close()
	j, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(j))
}

// GetSiteTemplates ...
func (s *Service) GetSiteTemplates(siteID string) {
	ep := fmt.Sprintf("/site/%v", siteID)
	res, err := s.http.MakeReq(s.baseURL+ep, "GET", nil)
	if err != nil {
	}
	defer res.Body.Close()
	j, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(j))
}

// Create ...
// func (s *Service) Create(name string) {
/*
	{
		"siteProfileUuid": "",
		"version" : 0,
		"name" name,
		"namespace" : "switching",
		"status": "",
		"lastUpdatedBy": "",
		"lastUpdatedDatetime": 0,
		"siteCount": 0,
		"profileAttributes": [],
		"id": "",
		"attributesList": [],
		"siteProfileType": "",
		"namingPrefix": "",
		"primaryDeviceType": "",
		"secondaryDeviceType": "",
		"interfaceList": [],
		"groupTypeList": [],
		"siteAssociationId": "",
		"sites": []
	}
*/
// }

// AssignSite ...
func (s *Service) AssignSite(profileID, siteID string) {
	// POST /siteprofile/{profileID}/site/{siteID}
}

// RemoveSite ...
func (s *Service) RemoveSite(profileID, siteID string) {
	/* DELETE /siteprofile/{profileID}/site/{siteID} */
}
