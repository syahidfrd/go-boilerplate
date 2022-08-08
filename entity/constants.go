package entity

type ctxKeyRequestID int

const RequestIDKey ctxKeyRequestID = 0

var RequestIDHeader = "X-Request-Id"
