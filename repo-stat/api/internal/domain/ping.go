package domain

type PingStatus string

const (
	PingStatusUp   PingStatus = "up"
	PingStatusDown PingStatus = "down"
)

type PingService struct {
	Name   string
	Status PingStatus
}

type PingResult struct {
	Status   string
	Services []PingService
}
