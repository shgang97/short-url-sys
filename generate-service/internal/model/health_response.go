package model

type HealthResponse struct {
	Status    string            `json:"status"`
	Datetime  string            `json:"datetime"`
	Timestamp int64             `json:"timestamp"`
	Services  map[string]string `json:"services,omitempty"`
}
