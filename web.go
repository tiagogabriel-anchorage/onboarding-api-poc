package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	NotFound    = Empty(http.StatusNotFound)
	ServerError = Empty(http.StatusInternalServerError)
)

type actionFunc func(r Requester) Responder

type Requester interface {
	ExtractBody(body any) error
	GetParamByName(name string) (string, bool)
}

type request struct {
	in *http.Request
}

func (r *request) ExtractBody(body any) error {
	defer r.in.Body.Close()
	if err := json.NewDecoder(r.in.Body).Decode(&body); err != nil {
		return err
	}
	return nil
}

func (r *request) GetParamByName(name string) (string, bool) {
	vars := mux.Vars(r.in)
	value, exists := vars[name]
	return value, exists
}

type Responder interface {
	WriteTo(out http.ResponseWriter)
}

type response struct {
	header http.Header
	status int
	body   []byte
}

func (r *response) WriteTo(out http.ResponseWriter) {
	header := out.Header()
	for k, v := range r.header {
		header[k] = v
	}
	out.WriteHeader(r.status)
	out.Write(r.body)
}

func (r *response) WithHeader(key, value string) *response {
	r.header.Set(key, value)
	return r
}

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

func Empty(status int) *response {
	return Content(status, struct{}{})
}

func EmptyError() *response {
	return ServerError
}

func Error(message string, err error) *response {
	return ErrorWithStatus(http.StatusInternalServerError, message, err)
}

func ErrorWithStatus(status int, message string, err error) *response {
	body := struct {
		message string
		details string
	}{
		message: message,
		details: err.Error(),
	}
	return newResponse(status, body).
		WithHeader("Content-Type", "application/json")
}

func Content(status int, body any) *response {
	return newResponse(status, body).
		WithHeader("Content-Type", "application/json")
}

func Ok(body any) *response {
	return Content(http.StatusOK, body)
}

func Accepted(body any) *response {
	return Content(http.StatusAccepted, body)
}

func Created(body any) *response {
	return Content(http.StatusCreated, body)
}

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
