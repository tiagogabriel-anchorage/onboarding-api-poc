package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func respondWithJson(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func extractedJsonRequest(r *http.Request, req any) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return ErrorResponse{
			Message:    "Something wrong with request",
			ErrMessage: err.Error(),
		}
	}
	return nil
}

func welcome(in Requester) Responder {
	return Ok(WelcomeResponse{Message: "Welcome to Onboarding API (PoC)"})
}

func postCustomers(in Requester) Responder {
	var req CreateCustomerRequest
	if err := in.ExtractBody(&req); err != nil {
		return BadRequest("Body is malformed", err)
	}

	// get kycSpec for this customer kind and entity
	kycSpec, err := getKYCSpec(req.CustomerKind, req.Entity)
	if err != nil {
		return BadRequest("Field values do not match", err)
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

	return Created(newCreateCustomerResponse(customer, kycSpec.KYCTemplate))
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
