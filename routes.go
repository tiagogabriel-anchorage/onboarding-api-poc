package main

import (
	"encoding/json"
	"errors"
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

func putCustomer(in Requester) Responder {
	var id uuid.UUID
	if value, exists := in.GetParamByName("id"); exists {
		var err error
		id, err = uuid.Parse(value)
		if err != nil {
			return BadRequest("Invalid id", fmt.Errorf("the id '%s' is an invalid universal identifier", value))
		}
	} else {
		return BadRequest("Missing route param", errors.New("'id' is missing"))
	}

	var req UpdateCustomerRequest
	if err := in.ExtractBody(&req); err != nil {
		return BadRequest("Body is malformed", err)
	}

	// get customer
	customer, exists := customerById(id)
	if !exists {
		return NotFound
	}

	if customer.KYCVersion != req.KYCVersion {
		return BadRequest("The KYC version is invalid",
			fmt.Errorf("kyc version '%d' is invalid", req.KYCVersion))
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

	return Ok(newUpdateCustomerResponse(customer))
}

func putCustomerV2(in Requester) Responder {
	var id uuid.UUID
	if value, exists := in.GetParamByName("id"); exists {
		var err error
		id, err = uuid.Parse(value)
		if err != nil {
			return BadRequest("Invalid id", fmt.Errorf("the id '%s' is an invalid universal identifier", value))
		}
	} else {
		return BadRequest("Missing route param", errors.New("'id' is missing"))
	}

	var req UpdateCustomerRequestV2
	if err := in.ExtractBody(&req); err != nil {
		return BadRequest("Body is malformed", err)
	}

	// get customer
	customer, exists := customerById(id)
	if !exists {
		return NotFound
	}

	if customer.KYCVersion != req.KYCVersion {
		return BadRequest("The KYC version is invalid",
			fmt.Errorf("kyc version '%d' is invalid", req.KYCVersion))
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

	return Ok(newUpdateCustomerResponse(customer))
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
