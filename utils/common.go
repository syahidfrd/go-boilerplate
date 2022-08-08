package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/syahidfrd/go-boilerplate/entity"
)

func CompactJSON(data []byte) string {
	var js map[string]interface{}
	if json.Unmarshal(data, &js) != nil {
		return string(data)
	}

	result := new(bytes.Buffer)
	if err := json.Compact(result, data); err != nil {
		fmt.Println(err)
	}
	return result.String()
}

// GetReqID get request id from echo context
func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(entity.RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}
