package trace

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey string

const traceKey ctxKey = "trace-id"

func NewContextWithTrace(ctx context.Context) context.Context {
	return context.WithValue(ctx, traceKey, uuid.New().String())
}

func GetTraceID(ctx context.Context) string {
	if v := ctx.Value(traceKey); v != nil {
		return v.(string)
	}
	return ""
}
