package entity

type ContextKey string

const (
	HeaderXCorrelationID  string     = "X-Correlation-ID"
	CorrelationContextKey ContextKey = "cid"
)
