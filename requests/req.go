package requests

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Req ...
type Req struct {
	http  *http.Client
	Token string
	user  string
	pass  string
}

// New creates instance of Req
func New(user, pass string) *Req {
	return &Req{
		http: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: 8 * time.Second,
		},
		user: user,
		pass: pass,
	}
}

// MakeReq ...
func (r *Req) MakeReq(uri, method string, b io.Reader) (*http.Response, error) {
	if !strings.Contains(uri, "/v1/auth") && r.Token == "" {
		return nil, fmt.Errorf("you must first login to perform this action")
	}
	req, err := http.NewRequest(method, uri, b)
	if err != nil {
		return nil, fmt.Errorf("unabled to create request: %v", err)
	}
	if r.Token != "" {
		req.Header.Set("X-AUTH-TOKEN", r.Token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}
	if strings.Contains(uri, "/v1/auth") {
		req.SetBasicAuth(r.user, r.pass)
	}
	res, err := r.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	if res.StatusCode == 401 {
		return nil, fmt.Errorf("token expired, %v", res.StatusCode)
	}
	return res, nil
}
