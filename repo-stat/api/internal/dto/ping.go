package dto

type PingResponse struct {
	Status   string            `json:"status"`
	Services []PingServiceInfo `json:"services"`
}

type PingServiceInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
