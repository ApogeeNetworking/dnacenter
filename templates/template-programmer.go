package templates

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

// New ...
func New(uri string, r *requests.Req) *Service {
	uri = strings.ReplaceAll(uri, "/dna", "")
	return &Service{baseURL: uri + "/api/v1/template-programmer", http: r}
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
	Name            string `json:"parameterName"`
	DataType        string `json:"dataType"`
	Description     string `json:"description"`
	Required        bool   `json:"required"`
	NotParam        bool   `json:"notParam"`
	DisplayName     string `json:"displayName"`
	InstructionText string `json:"instructionText"`
	Order           int    `json:"order"`
	Binding         string `json:"binding"`
	ID              string `json:"id"`
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

// PreviewCfgResp ...
type PreviewCfgResp struct {
	TemplateID string `json:"templateId"`
	CliPreview string `json:"cliPreview"`
}

// GenPreviewCfg ...
func (s *Service) GenPreviewCfg(v interface{}) (PreviewCfgResp, error) {
	/*
		{
			params: {
				ParamKey: ParamKeyValue
				...
			},
			templateId: templateId
		}
	*/
	var data []byte
	switch v.(type) {
	case []byte:
		data = v.([]byte)
	default:
		data, _ = json.Marshal(&v)
	}
	body := strings.NewReader(string(data))
	endpoint := "/template/preview"
	res, err := s.http.MakeReq(s.baseURL+endpoint, "PUT", body)
	if err != nil {
		return PreviewCfgResp{}, err
	}
	defer res.Body.Close()
	var preview PreviewCfgResp
	json.NewDecoder(res.Body).Decode(&preview)
	return preview, nil
}
