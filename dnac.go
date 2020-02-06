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

// Client used for DNA-C Connection Handler
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

// MakeReq on behalf of our DNAC Client
func (c *Client) MakeReq(path, method string, r io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.BaseURL+path, r)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	if c.authToken != "" {
		req.Header.Set("X-AUTH-TOKEN", c.authToken)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}
	if strings.Contains(path, "/v1/auth") {
		req.SetBasicAuth(c.Username, c.Password)
	}
	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	return res, nil
}

// Login establishes a session with the DNAC API
func (c *Client) Login() error {
	res, err := c.MakeReq("/api/system/v1/auth/token", "POST", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var tok struct {
		Token string `json:"Token"`
	}
	json.NewDecoder(res.Body).Decode(&tok)
	// Retain AuthToken for Client (keeping it private)
	c.authToken = tok.Token
	// Reset the BaseURL
	c.BaseURL = fmt.Sprintf("https://%s/dna", c.IP)

	return nil
}
