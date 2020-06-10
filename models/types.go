package models

// Task ...
type Task struct {
	TaskID      string `json:"taskId,omitempty"`
	ID          string `json:"id,omitempty"`
	Progess     string `json:"progress,omitempty"`
	IsError     bool   `json:"isError,omitempty"`
	ServiceType string `json:"serviceType,omitempty"`
	User        string `json:"username,omitempty"`
	URL         string `json:"url,omitempty"`
}
