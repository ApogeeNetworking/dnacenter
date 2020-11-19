package events

import (
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
func New(uri string, c *requests.Req) *Service {
	uri = strings.ReplaceAll(uri, "/dna", "")
	return &Service{
		baseURL: fmt.Sprintf("%s/api/v1/event", uri),
		http:    c,
	}
}

// GetSubscriptions ...
func (s *Service) GetSubscriptions() error {
	_, err := s.http.MakeReq(s.baseURL+"/subscription", "GET", nil)
	if err != nil {
		return err
	}
	return nil
}
