package templates

import (
	"encoding/json"
	"fmt"

	"github.com/ApogeeNetworking/dnacenter/requests"
)

// Service ...
type Service struct {
	baseURL string
	http    *requests.Req
}

// New ...
func New(uri string, r *requests.Req) *Service {
	return &Service{baseURL: uri + "/intent/api/v1/template-programmer", http: r}
}

// Project a Simple DNAC Project
type Project struct {
	Name        string     `json:"name"`
	ID          string     `json:"id"`
	Templates   []Template `json:"templates"`
	IsDeletable bool       `json:"isDeletable"`
	Description string     `json:"description,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	Created     int64      `json:"createTime,omitempty"`
	LastUpdated int64      `json:"lastUpdateTime,omitempty"`
}

// GetProjects retrieves projects and templates assoc w/ea
// func (c *Client) GetProjects() ([]Project, error) {
// 	if c.authToken == "" {
// 		return nil, fmt.Errorf(loginWarning)
// 	}
// 	ep := "/intent/api/v1/template-programmer/project"
// 	res, err := c.MakeReq(ep, "GET", nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()
// 	var projects []Project
// 	json.NewDecoder(res.Body).Decode(&projects)
// 	return projects, nil
// }

// CreateProject adds a project in DNAC
// **DO NOT USE** Apparently Projects Created via the API
// are NON-DELETABLE in DNA-C
// func (c *Client) CreateProject(p Project) (DNARes, error) {
// 	if c.authToken == "" {
// 		return DNARes{}, fmt.Errorf(loginWarning)
// 	}
// 	if p.Name == "" {
// 		return DNARes{}, fmt.Errorf("bad request")
// 	}
// 	j, _ := json.Marshal(p)
// 	body := strings.NewReader(string(j))
// 	ep := "/intent/api/v1/template-programmer/project"
// 	res, err := c.MakeReq(ep, "POST", body)
// 	if err != nil {
// 		return DNARes{}, err
// 	}
// 	defer res.Body.Close()
// 	r, _ := ioutil.ReadAll(res.Body)
// 	fmt.Println(string(r))
// 	var result DNARes
// 	json.NewDecoder(res.Body).Decode(&result)
// 	return result, nil
// }

// DeleteProject deletes a project in DNAC by projectID
// func (c *Client) DeleteProject(projectID string) (DNARes, error) {
// 	if c.authToken == "" {
// 		return DNARes{}, fmt.Errorf(loginWarning)
// 	}
// 	ep := fmt.Sprintf("/intent/api/v1/template-programmer/project/%s", projectID)
// 	res, err := c.MakeReq(ep, "DELETE", nil)
// 	if err != nil {
// 		return DNARes{}, err
// 	}
// 	defer res.Body.Close()
// 	fmt.Println(res.StatusCode)
// 	var result DNARes
// 	json.NewDecoder(res.Body).Decode(&result)
// 	return result, nil
// }

// Template a DNAC PnP Config "File"
type Template struct {
	Name               string            `json:"name"`
	Composite          bool              `json:"composite"`
	ID                 string            `json:"id"`
	Description        string            `json:"description,omitempty"`
	Tags               []string          `json:"tag,omitempty"`
	DeviceTypes        []TemplDeviceType `json:"deviceTypes,omitempty"`
	SoftwareType       string            `json:"softwareType,omitempty"`
	SoftwareVariant    string            `json:"softwareVariant,omitempty"`
	Content            string            `json:"templateContent,omitempty"`
	Params             []TemplParam      `json:"templateParams,omitempty"`
	RollbackContent    string            `json:"rollbackTemplateContent,omitempty"`
	RollbackParams     []string          `json:"rollbackTemplateParams,omitempty"`
	ContainingTemplate []string          `json:"containingTemplate,omitempty"`
	Created            int64             `json:"createTime,omitempty"`
	LastUpdated        int64             `json:"lastUpdateTime,omitempty"`
	ProjectName        string            `json:"projectName,omitempty"`
	ProjectID          string            `json:"projectId,omitempty"`
	ParentID           string            `json:"parentTemplateId,omitempty"`
}

// TemplDeviceType for a DNAC Template
type TemplDeviceType struct {
	ProductFamily string `json:"productFamily"`
	ProductSeries string `json:"productSeries"`
	ProductType   string `json:"productType"`
}

// TemplParam $Var delimiters in a Template
type TemplParam struct {
	Name            string        `json:"parameterName"`
	DataType        string        `json:"dataType"`
	DefaultValue    interface{}   `json:"defaultValue"`
	Description     string        `json:"description"`
	Required        bool          `json:"required"`
	NotParam        bool          `json:"notParam"`
	DisplayName     string        `json:"displayName"`
	InstructionText string        `json:"instructionText"`
	Group           interface{}   `json:"group"`
	Order           int           `json:"order"`
	Selection       interface{}   `json:"selection"`
	Range           []interface{} `json:"range"`
	Key             interface{}   `json:"key"`
	Provider        interface{}   `json:"provider"`
	Binding         string        `json:"binding"`
	ID              string        `json:"id"`
}

// GetTemplate retrieves a Template from DNAC by ID
func (s *Service) GetTemplate(templateID string) (Template, error) {
	ep := fmt.Sprintf("/template/%s", templateID)
	res, err := s.http.MakeReq(s.baseURL+ep, "GET", nil)
	if err != nil {
		return Template{}, err
	}
	defer res.Body.Close()
	var template Template
	json.NewDecoder(res.Body).Decode(&template)
	return template, nil
}
