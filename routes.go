package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func respondWithJson(w http.ResponseWriter, statusCode int, body any) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Could not parse the response. Err: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
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
	now := time.Now()
	customer := Customer{
		ID:         uuid.MustParse("878bc7da-4809-11ed-b878-0242ac120002"),
		KYCVersion: kycSpec.KYCVersion,
		Status:     "DRAFT",
		CreatedAt:  now,
		UpdatedAt:  now,
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
	customer, exists := customerById(id)
	if !exists {
		respondWithJson(w, http.StatusNotFound, ErrorResponse{
			Message: fmt.Sprintf("The customer '%s' does not exist", id),
		})
		return
	}

	if customer.KYCVersion != req.KYCVersion {
		respondWithJson(w, http.StatusBadRequest, ErrorResponse{
			Message: "The KYC version is invalid",
		})
		return
	}

	/**
	* TODO:
	* (1) validate if field is valid for given answer
	* (2) match the exposed question_ids to the internal ones
	 */
	var answers []Answer
	for _, kycEntry := range req.KYC {
		answers = append(answers, Answer{
			QuestionID: kycEntry.QuestionID,
			Answer:     kycEntry.Answer,
		})
	}

	customer.SubmitAnswers(answers)
	customer.UpdatedAt = time.Now()
	saveKYC(customer)

	respondWithJson(w, http.StatusOK, newUpdateCustomerResponse(customer))
}

func putCustomerV2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.MustParse(vars["id"])

	var req UpdateCustomerRequestV2
	if err := extractedJsonRequest(r, &req); err != nil {
		respondWithJson(w, http.StatusBadRequest, err)
		return
	}

	// get customer
	customer, exists := customerById(id)
	if !exists {
		respondWithJson(w, http.StatusNotFound, ErrorResponse{
			Message: fmt.Sprintf("The customer '%s' does not exist", id),
		})
		return
	}

	if customer.KYCVersion != req.KYCVersion {
		respondWithJson(w, http.StatusBadRequest, ErrorResponse{
			Message: "The KYC version is invalid",
		})
		return
	}

	/**
	* TODO:
	* (1) validate if field is valid for given answer
	* (2) match the exposed question_ids to the internal ones
	 */
	var answers []Answer
	for _, kycEntry := range req.KYC {
		sepIndex := strings.Index(kycEntry, ":")
		answers = append(answers, Answer{
			QuestionID: kycEntry[:sepIndex],
			Answer:     kycEntry[sepIndex+1:], // this answer should be parsed according to the answer type
		})
	}

	customer.SubmitAnswers(answers)
	customer.UpdatedAt = time.Now()
	saveKYC(customer)

	respondWithJson(w, http.StatusOK, newUpdateCustomerResponse(customer))
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.MustParse(vars["id"])

	// get customer
	customer, exists := customerById(id)
	if !exists {
		respondWithJson(w, http.StatusNotFound, ErrorResponse{
			Message: fmt.Sprintf("The customer '%s' does not exist", id),
		})
		return
	}

	respondWithJson(w, http.StatusOK, newGetCustomerResponse(customer))
}

func postCustomerSubmission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := uuid.MustParse(vars["id"])

	// get customer
	customer, exists := customerById(id)
	if !exists {
		respondWithJson(w, http.StatusNotFound, ErrorResponse{
			Message: fmt.Sprintf("The customer '%s' does not exist", id),
		})
		return
	}

	customer.Status = "LATEST"
	saveKYC(customer)

	respondWithJson(w, http.StatusOK, newGetCustomerResponse(customer))
}
