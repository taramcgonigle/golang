package api

import (
	"net/http"
)

// RegisterRoutes sets up all the API routes
func RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", withTraceID(http.HandlerFunc(HomeHandler)))
	mux.Handle("/create-form", withTraceID(http.HandlerFunc(CreateFormHandler)))
	mux.Handle("/create", withTraceID(http.HandlerFunc(CreateHandler)))
	mux.Handle("/get", withTraceID(http.HandlerFunc(GetHandler)))
	mux.Handle("/edit-form", withTraceID(http.HandlerFunc(EditFormHandler)))
	mux.Handle("/update", withTraceID(http.HandlerFunc(UpdateHandler)))
	mux.Handle("/delete", withTraceID(http.HandlerFunc(DeleteHandler)))

	return mux
}
