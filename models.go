package main

import (
	"time"

	"github.com/google/uuid"
)

// WelcomeResponse is the message of welcoming.
type WelcomeResponse struct {
	Message string `json:"message"`
}

// CreateCustomerRequest represents the payload for creating new customers.
type CreateCustomerRequest struct {
	CustomerKind string `json:"kind"`
	Entity       string `json:"entity"`
}

// CreateCustomerResponse represents the response payload of a new customer creation.
type CreateCustomerResponse struct {
	ID          uuid.UUID             `json:"customer_id"`
	KYCVersion  int                   `json:"kyc_version"`
	KYCTemplate []KYCTemplateResponse `json:"kyc_template"`
	Status      string                `json:"status"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

// KYCTemplateResponse represents the node of the KYC specification.
type KYCTemplateResponse struct {
	QuestionID string   `json:"question_id"`
	AnswerType string   `json:"answer_type"`
	DependsOn  []string `json:"depends_on"`
	Mandatory  bool     `json:"mandatory"`
}

// UpdateCustomerRequest represents the payload for updating KYC information of a given customer.
type UpdateCustomerRequest struct {
	KYCVersion int `json:"kyc_version"`
	KYC        []struct {
		QuestionID string `json:"question_id"`
		Answer     any    `json:"answer"`
	} `json:"kyc"`
}

// UpdateCustomerRequestV2 represents the payload for updating KYC information of a given customer,
// but with a different representation (bounded to a new version endpoint).
type UpdateCustomerRequestV2 struct {
	KYCVersion int      `json:"kyc_version"`
	KYC        []string `json:"kyc"`
}

// UpdateCustomerResponse represents the response payload of a KYC information update, for a given customer.
type UpdateCustomerResponse struct {
	ID         uuid.UUID `json:"customer_id"`
	KYCVersion int       `json:"kyc_version"`
	KYC        []struct {
		QuestionID string `json:"question_id"`
		Answer     any    `json:"answer"`
	} `json:"kyc"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetCustomerResponse represents the current state of KYC information for a given customer.
type GetCustomerResponse struct {
	ID         uuid.UUID `json:"customer_id"`
	KYCVersion int       `json:"kyc_version"`
	KYC        []struct {
		QuestionID string `json:"question_id"`
		Answer     any    `json:"answer"`
	} `json:"kyc"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// newCreateCustomerResponse maps a customer created to the response to be delivered to the client.
func newCreateCustomerResponse(customer Customer, kycTemplate []KYCTemplateEntry) CreateCustomerResponse {
	res := CreateCustomerResponse{
		ID:         customer.ID,
		KYCVersion: customer.KYCVersion,
		Status:     customer.Status,
		CreatedAt:  customer.CreatedAt,
		UpdatedAt:  customer.UpdatedAt,
	}

	for _, entry := range kycTemplate {
		res.KYCTemplate = append(res.KYCTemplate, KYCTemplateResponse{
			QuestionID: entry.QuestionID,
			AnswerType: entry.AnswerType,
			DependsOn:  entry.DependsOn,
			Mandatory:  entry.Mandatory,
		})
	}

	return res
}

// newUpdateCustomerResponse maps an updated KYC information of a given customer to the response to be delivered to the client.
func newUpdateCustomerResponse(customer Customer) UpdateCustomerResponse {
	res := UpdateCustomerResponse{
		ID:         customer.ID,
		KYCVersion: customer.KYCVersion,
		Status:     customer.Status,
		CreatedAt:  customer.CreatedAt,
		UpdatedAt:  customer.UpdatedAt,
	}

	for _, entry := range customer.Answers {
		res.KYC = append(res.KYC, struct {
			QuestionID string `json:"question_id"`
			Answer     any    `json:"answer"`
		}{
			QuestionID: entry.QuestionID,
			Answer:     entry.Answer,
		})
	}

	return res
}

// newGetCustomerResponse maps the current state of KYC information of a given customer to the response to be delivered to the client.
func newGetCustomerResponse(customer Customer) GetCustomerResponse {
	res := GetCustomerResponse{
		ID:         customer.ID,
		KYCVersion: customer.KYCVersion,
		Status:     customer.Status,
		KYC: []struct {
			QuestionID string "json:\"question_id\""
			Answer     any    "json:\"answer\""
		}{},
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}

	for _, kycEntry := range customer.Answers {
		res.KYC = append(res.KYC, struct {
			QuestionID string "json:\"question_id\""
			Answer     any    "json:\"answer\""
		}{
			QuestionID: kycEntry.QuestionID,
			Answer:     kycEntry.Answer,
		})
	}

	return res
}
