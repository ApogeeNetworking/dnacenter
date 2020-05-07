package sites

import (
	"encoding/json"

	"github.com/drkchiloll/dnacenter/requests"
)

// Service ...
type Service struct {
	baseURL string
	http    *requests.Req
}

// New creates an instance of a DNA-C Sites Service
func New(uri string, r *requests.Req) *Service {
	return &Service{
		baseURL: uri + "/intent/api/v1/site",
		http:    r,
	}
}

// Site properties of a Site in DNA-C
type Site struct {
	ParentID          string   `json:"parentId,omitempty"`
	Name              string   `json:"name"`
	AdditionalInfo    []string `json:"additionalInfo"`
	SiteHierarchy     string   `json:"siteHierarchy"`
	SiteNameHierarchy string   `json:"siteNameHierarchy"`
	InstanceTenantID  string   `json:"instanceTenantId"`
	ID                string   `json:"id"`
}

// SiteParams ...
type SiteParams struct {
	// (Optional) Fully Qualified Site Name
	// ex: Global/<SiteName>
	Name string
	// (Optional) SiteID
	SiteID string
	// (Optional) Type = area, building, floor
	Type string
	// (Optional) OffSet = starting row
	OffSet int
	// (Optional) Limit = number of sites to be retrieved
	Limit int
}

// Get retrieves all sites
func (s *Service) Get(p SiteParams) ([]Site, error) {
	uri := s.baseURL
	if p.Name != "" {
		uri += "?name=" + p.Name
	}
	res, err := s.http.MakeReq(uri, "GET", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	type resp struct {
		Response []Site `json:"response"`
	}
	var sites resp
	json.NewDecoder(res.Body).Decode(&sites)
	return sites.Response, err
}
