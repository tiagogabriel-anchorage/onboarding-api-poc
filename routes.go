package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	var req CreateCustomerRequest
	if err := extractedJsonRequest(r, &req); err != nil {
		respondWithJson(w, http.StatusBadRequest, err)
		return
	}

	// get kycSpec for this customer kind and entity
	kycSpec, err := getKYCSpec(req.CustomerKind, req.Entity)
	if err != nil {
		respondWithJson(w, http.StatusBadRequest, ErrorResponse{
			Message:    "Field values do not match",
			ErrMessage: err.Error(),
		})
		return
	}

	// if spec available, create the customer and save
	customer := Customer{
		CustomerID: uuid.MustParse("878bc7da-4809-11ed-b878-0242ac120002"),
		KYCVersion: kycSpec.KYCVersion,
		Status:     "DRAFT",
	}
	saveKYC(customer)

	// fill response from kyc spec
	respondWithJson(w, http.StatusCreated, newCreateCustomerResponse(customer, kycSpec.KYCTemplate))
}

func putCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.MustParse(vars["id"])

	var req UpdateCustomerRequest
	if err := extractedJsonRequest(r, &req); err != nil {
		respondWithJson(w, http.StatusBadRequest, err)
		return
	}

	// get customer
	customer := customerById(id)
	if customer.KYCVersion != req.KYCVersion {
		respondWithJson(w, http.StatusBadRequest, ErrorResponse{
			Message: "The KYC version is invalid",
		})
		return
	}

	// fill in answers
	customer.Answers = nil // forget about previous answers
	for _, kycEntry := range req.KYC {
		customer.Answers = append(customer.Answers, Answer{
			QuestionID: kycEntry.QuestionID,
			Answer:     kycEntry.Answer,
		})
	}
	// Persist the customer
	saveKYC(customer)

	// Respond
	respondWithJson(w, http.StatusOK, newUpdateCustomerResponse(customer))
}
