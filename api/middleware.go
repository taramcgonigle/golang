package api

import (
	"log/slog"
	"net/http"
	"practice/assignments/trace"
)

// Middleware to inject a trace ID into the request context
func withTraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := trace.NewContextWithTrace(r.Context())
		traceID := trace.GetTraceID(ctx)

		slog.Info("Incoming request", "traceID", traceID, "method", r.Method, "path", r.URL.Path)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
