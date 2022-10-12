package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Welcome responds to the request to the root path.
func Welcome(in Requester) Responder {
	return Ok(WelcomeResponse{Message: "Welcome to Onboarding API (PoC)"})
}

// PostCustomer deals with the registry of a new customer KYC information.
func PostCustomers(in Requester) Responder {
	var req CreateCustomerRequest
	if err := in.ExtractBody(&req); err != nil {
		return BadRequest("Body is malformed", err)
	}

	// get kycSpec for this customer kind and entity
	kycSpec, err := getKYCSpec(req.CustomerKind, req.Entity)
	if err != nil {
		return BadRequest("Could not find a KYC specification", err)
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

// PutCustomer submits the answers to the current KYC information of the customer.
func PutCustomer(in Requester) Responder {
	id, res := getCustomerID(in)
	if res != nil {
		return res
	}

	var req UpdateCustomerRequest
	if err := in.ExtractBody(&req); err != nil {
		return BadRequest("Body is malformed", err)
	}

	// get customer
	customer, exists := customerByID(id)
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

// PutCustomer submits the answers to the current KYC information of the customer,
// but uses a different approach when dealing with the input data.
func PutCustomerV2(in Requester) Responder {
	id, res := getCustomerID(in)
	if res != nil {
		return res
	}

	var req UpdateCustomerRequestV2
	if err := in.ExtractBody(&req); err != nil {
		return BadRequest("Body is malformed", err)
	}

	// get customer
	customer, exists := customerByID(id)
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
			Answer:     kycEntry[sepIndex+1:], // this answer should be parsed accordingly to the answer type
		})
	}

	customer.SubmitAnswers(answers)
	customer.UpdatedAt = time.Now()
	saveKYC(customer)

	return Ok(newUpdateCustomerResponse(customer))
}

// GetCustomer retrieves the current state of the KYC for the given customer.
func GetCustomer(in Requester) Responder {
	id, res := getCustomerID(in)
	if res != nil {
		return res
	}

	// get customer
	customer, exists := customerByID(id)
	if !exists {
		return NotFound
	}

	return Ok(newGetCustomerResponse(customer))
}

// PostCustomerSubmission marks the KYC information as final, therefore submitting the KYC information.
func PostCustomerSubmission(in Requester) Responder {
	id, res := getCustomerID(in)
	if res != nil {
		return res
	}

	// get customer
	customer, exists := customerByID(id)
	if !exists {
		return NotFound
	}

	customer.Status = "LATEST"
	saveKYC(customer)

	return Ok(newGetCustomerResponse(customer))
}

// getCustomerID extracts the id value from the URL path.
func getCustomerID(in Requester) (uuid.UUID, Responder) {
	var id uuid.UUID
	if value, exists := in.GetParamByName("id"); exists {
		var err error
		id, err = uuid.Parse(value)
		if err != nil {
			return id, BadRequest("Invalid id", fmt.Errorf("the id '%s' is an invalid universal identifier", value))
		}
	} else {
		return id, BadRequest("Missing route param", errors.New("'id' is missing"))
	}
	return id, nil
}
