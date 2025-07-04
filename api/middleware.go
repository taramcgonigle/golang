package api

import (
	"log/slog"
	"net/http"
	"practice/assignments/trace"
)

// Middleware to inject a trace ID into the request context for logging
func withTraceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //accepts any http.Handler and reutns a new one that adds extra behavior
		ctx := trace.NewContextWithTrace(r.Context()) //generate a new trace ID and add it to the request context
		traceID := trace.GetTraceID(ctx)

		slog.Info("Incoming request", "traceID", traceID, "method", r.Method, "path", r.URL.Path) //logs request

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
