package dnac

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const loginWarning string = "you must first login to perform this action"

// Client used for CDNAC Connection Handler
type Client struct {
	BaseURL  string
	Username string
	Password string
	IP       string
	// These are purposely not exported
	http      *http.Client
	authToken string
}

// NewClient creates a reference to the DNAC Client Struct
func NewClient(host, user, pass string, ignoreSSL bool) *Client {
	return &Client{
		BaseURL:  fmt.Sprintf("https://%s", host),
		Username: user,
		Password: pass,
		IP:       host,
		http: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: ignoreSSL,
				},
			},
			Timeout: 8 * time.Second,
		},
	}
}

func (c *Client) genReq(url, method string, r io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.BaseURL+url, r)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	if c.authToken != "" {
		req.Header.Set("X-AUTH-TOKEN", c.authToken)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}
	if strings.Contains(url, "/v1/auth") {
		req.SetBasicAuth(c.Username, c.Password)
	}

	return req, nil
}

// Login establishes a session with the DNAC API
func (c *Client) Login() error {
	// First we have to perform an Auth Login
	req, err := c.genReq("/api/system/v1/auth/login", "GET", nil)
	if err != nil {
		return err
	}
	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	// Grab the Cookie for the Next REQ
	// We Do Not Have to Do this Afterward
	authCookie := res.Cookies()[0]

	c.BaseURL = fmt.Sprintf("https://%s/dna", c.IP)

	req, err = c.genReq("/system/api/v1/auth/token", "POST", nil)
	if err != nil {
		return fmt.Errorf("failed to create a request: %v", err)
	}
	req.AddCookie(authCookie)

	res, err = c.http.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer res.Body.Close()
	// var tok AuthTok
	var tok struct {
		Token string `json:"Token"`
	}
	json.NewDecoder(res.Body).Decode(&tok)
	c.authToken = tok.Token
	return nil
}