package dnac

import (
	"encoding/json"
	"fmt"

	"github.com/ApogeeNetworking/dnacenter/devices"
	"github.com/ApogeeNetworking/dnacenter/pnp"
	"github.com/ApogeeNetworking/dnacenter/requests"
	siteprofile "github.com/ApogeeNetworking/dnacenter/site-profile"
	"github.com/ApogeeNetworking/dnacenter/sites"
	"github.com/ApogeeNetworking/dnacenter/swim"
	"github.com/ApogeeNetworking/dnacenter/templates"
)

const loginWarning string = "you must first login to perform this action"

// Client used for DNA-C Connection Handler
type Client struct {
	BaseURL string
	IP      string
	// These are purposely not exported
	http        *requests.Req
	Devices     *devices.Service
	Sites       *sites.Service
	SWIM        *swim.Service
	SiteProfile *siteprofile.Service
	PnP         *pnp.Service
	Templates   *templates.Service
}

// NewClient creates a reference to the DNAC Client Struct
func NewClient(host, user, pass string, ignoreSSL bool) *Client {
	client := &Client{
		BaseURL: fmt.Sprintf("https://%s/dna", host),
		IP:      host,
		http:    requests.New(user, pass),
	}
	client.Sites = sites.New(client.BaseURL, client.http)
	client.Devices = devices.New(client.BaseURL, client.http)
	client.SWIM = swim.New(client.BaseURL, client.http)
	client.SiteProfile = siteprofile.New(client.BaseURL, client.http)
	client.PnP = pnp.New(client.BaseURL, client.http)
	client.Templates = templates.New(client.BaseURL, client.http)
	return client
}

// Login establishes a session with the DNAC API
func (c *Client) Login() error {
	authURL := "/system/api/v1/auth/token"
	res, err := c.http.MakeReq(c.BaseURL+authURL, "POST", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var tok struct {
		Token string `json:"Token"`
	}
	json.NewDecoder(res.Body).Decode(&tok)
	// Retain AuthToken for REQ Client (keeping it private)
	c.http.Token = tok.Token
	return nil
}

// Version in most DNA-C Responses
type Version string

// DNARes object returned from POST/PUT/DELETE REQS
type DNARes struct {
	Response struct {
		TaskID string `json:"taskId"`
		URL    string `json:"url"`
	} `json:"response"`
	Version `json:"version"`
}

// DNATaskCheckRes properties when checking a Tasks Progress
type DNATaskCheckRes struct {
	End         int64  `json:"endTime"`
	Version     int64  `json:"version"`
	Progress    string `json:"progress"`
	Start       int64  `json:"startTime"`
	Username    string `json:"username"`
	Data        string `json:"data"`
	IsError     bool   `json:"isError"`
	RootID      string `json:"rootId"`
	ServiceType string `json:"serviceType"`
	InstanceID  string `json:"instanceId"`
	ID          string `json:"id"`
}

// DNATaskRes results of a Task Request
type DNATaskRes struct {
	Response DNATaskCheckRes `json:"response"`
	Version  `json:"version"`
}
