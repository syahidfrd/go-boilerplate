package health

// Status represents the health status of a component
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
)

// HealthResponse represents the overall health check response
type HealthResponse struct {
	Database Status `json:"database"`
	Cache    Status `json:"cache"`
}
