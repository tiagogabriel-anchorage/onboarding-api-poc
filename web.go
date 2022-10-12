package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	// NotFound defines a resource that does not exist.
	NotFound = Empty(http.StatusNotFound)
	// ServerError specifies an unspecified error in the server.
	ServerError = Empty(http.StatusInternalServerError)
)

// A actionFunc is a function that specifies the behavior of an endpoint.
type actionFunc func(r Requester) Responder

// A Requester allows interaction with the information sent by the client.
type Requester interface {
	ExtractBody(body any) error
	GetParamByName(name string) (string, bool)
}

// A request wraps the original HTTP request.
type request struct {
	in *http.Request
}

// ExtractBody decodes the body of a HTTP request into a data structure.
func (r *request) ExtractBody(body any) error {
	defer r.in.Body.Close()
	if err := json.NewDecoder(r.in.Body).Decode(&body); err != nil {
		return err
	}
	return nil
}

// GetParamByName retrieves a param value by its name (includes route, query params).
func (r *request) GetParamByName(name string) (string, bool) {
	vars := mux.Vars(r.in)
	value, exists := vars[name]
	return value, exists
}

// A Responder proceeds with the communication with the client.
type Responder interface {
	WriteTo(out http.ResponseWriter)
}

// response wraps the data to be deliver to the client.
type response struct {
	header http.Header
	status int
	body   []byte
}

// WriteTo writes the response produced by the server.
func (r *response) WriteTo(out http.ResponseWriter) {
	header := out.Header()
	for k, v := range r.header {
		header[k] = v
	}
	out.WriteHeader(r.status)
	out.Write(r.body)
}

// WithHeader specifies headers for the response payload.
func (r *response) WithHeader(key, value string) *response {
	r.header.Set(key, value)
	return r
}

// newResponse produces a new response, given a status and a body.
func newResponse(status int, body any) *response {
	var content []byte
	var err error
	switch t := body.(type) {
	case []byte:
		content = t
	case string:
		content = []byte(t)
	default:
		if content, err = json.Marshal(body); err != nil {
			return ErrorWithStatus(http.StatusBadRequest, "Could not parse body", err)
		}
	}
	return &response{
		header: http.Header{},
		status: status,
		body:   content,
	}
}

// Empty provides empty response.
func Empty(status int) *response {
	return Content(status, struct{}{})
}

// EmptyError provides a ServerError response without body.
func EmptyError() *response {
	return ServerError
}

// Error provides a ServerError with body containing relevant info.
func Error(message string, err error) *response {
	return ErrorWithStatus(http.StatusInternalServerError, message, err)
}

// ErrorWithStatus provides a body containing relevant info, and overriding the status.
func ErrorWithStatus(status int, message string, err error) *response {
	body := struct {
		Message string `json:"message"`
		Details string `json:"details"`
	}{
		Message: message,
		Details: err.Error(),
	}
	return newResponse(status, body).
		WithHeader("Content-Type", "application/json")
}

// Content provides a content to the client.
func Content(status int, body any) *response {
	return newResponse(status, body).
		WithHeader("Content-Type", "application/json")
}

// Ok provides content to the client, with the 200 OK status.
func Ok(body any) *response {
	return Content(http.StatusOK, body)
}

// Created provides content to the client, with the 201 Accepted status.
func Created(body any) *response {
	return Content(http.StatusCreated, body)
}

// Accepted provides content to the client, with the 202 Accepted status.
func Accepted(body any) *response {
	return Content(http.StatusAccepted, body)
}

// BadRequest provides a body containing relevant info, with the 400 BadRequest status.
func BadRequest(message string, err error) *response {
	return ErrorWithStatus(http.StatusBadRequest, message, err)
}

// handleRoute configures the router with path, action to execute and verbs associated with the endpoint.
func handleRoute(r *mux.Router, path string, action actionFunc, httpVerbs ...string) {
	handler := func(out http.ResponseWriter, in *http.Request) {
		res := action(&request{in: in})
		if res == nil {
			res = ServerError
		}
		res.WriteTo(out)
	}
	r.HandleFunc(path, handler).
		Methods(httpVerbs...)
}
