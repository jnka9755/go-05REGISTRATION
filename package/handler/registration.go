package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/digitalhouse-tech/go-lib-kit/response"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jnka9755/go-05REGISTRATION/internal/registration"
)

func NewCourseHTTPServer(ctx context.Context, endpoints registration.Endpoints) http.Handler {

	r := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Handle("/registration", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateRegistration,
		encodeResponse,
		opts...,
	)).Methods("POST")

	return r
}

func decodeCreateRegistration(_ context.Context, r *http.Request) (interface{}, error) {

	var req registration.CreateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("Invalid reques format: '%v'", err.Error()))
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	r := err.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	_ = json.NewEncoder(w).Encode(r)
}
