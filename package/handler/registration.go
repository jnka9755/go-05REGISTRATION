package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jnka9755/go-05REGISTRATION/internal/registration"
	"github.com/jnka9755/go-05RESPONSE/response"
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

	r.Handle("/registration", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllRegistration,
		encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/registration/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateRegistration,
		encodeResponse,
		opts...,
	)).Methods("PATCH")

	return r
}

func decodeCreateRegistration(_ context.Context, r *http.Request) (interface{}, error) {

	var req registration.CreateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("Invalid reques format: '%v'", err.Error()))
	}

	return req, nil
}

func decodeGetAllRegistration(_ context.Context, r *http.Request) (interface{}, error) {

	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := registration.GetAllReq{
		UserID:   v.Get("user_id"),
		CourseID: v.Get("course_id"),
		Limit:    limit,
		Page:     page,
	}

	return req, nil
}

func decodeUpdateRegistration(_ context.Context, r *http.Request) (interface{}, error) {

	var req registration.UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("Invalid request format: '%v'", err.Error()))
	}

	p := mux.Vars(r)
	req.ID = p["id"]

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
