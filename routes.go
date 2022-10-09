package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func respondWithJson(w http.ResponseWriter, statusCode int, body any) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Could not parse the response. Err: %s", err)
	}
	w.WriteHeader(statusCode)
	w.Write(jsonBody)
}

func extractedJsonRequest(r *http.Request, req any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return ErrorResponse{
			Message:    "Something wrong parsing the body.",
			ErrMessage: err.Error(),
		}
	}

	if err := json.Unmarshal(body, &req); err != nil {
		return ErrorResponse{
			Message:    "Something wrong unmarshalling body",
			ErrMessage: err.Error(),
		}
	}

	return nil
}

func welcome(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusOK, WelcomeResponse{Message: "Welcome to Onboarding API (PoC)"})
}

func postCustomers(w http.ResponseWriter, r *http.Request) {
	var req NewCustomerRequest
	if err := extractedJsonRequest(r, &req); err != nil {
		respondWithJson(w, http.StatusBadRequest, err)
		return
	}

	// Lets support, for now, only business and anchorage hold entity
	if !strings.EqualFold(req.CustomerKind, "business") {
		respondWithJson(w, http.StatusBadRequest, ErrorResponse{
			Message:    "Does not support the given customer type",
			ErrMessage: fmt.Sprintf("'%s' not supported", req.CustomerKind),
		})
		return
	}

	if !strings.EqualFold(req.Entity, "anchorage hold") {
		respondWithJson(w, http.StatusBadRequest, ErrorResponse{
			Message:    "Does not support the given entity",
			ErrMessage: fmt.Sprintf("'%s' not supported for customer type", req.Entity),
		})
		return
	}

	respondWithJson(w, http.StatusCreated, req)
}
