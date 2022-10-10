package main

import (
	"time"

	"github.com/google/uuid"
)

type ErrorResponse struct {
	error
	Message    string `json:"message"`
	ErrMessage string `json:"error"`
}

type WelcomeResponse struct {
	Message string `json:"message"`
}

type CreateCustomerRequest struct {
	CustomerKind string `json:"kind"`
	Entity       string `json:"entity"`
}

type CreateCustomerResponse struct {
	ID          uuid.UUID             `json:"customer_id"`
	KYCVersion  int                   `json:"kyc_version"`
	KYCTemplate []KYCTemplateResponse `json:"kyc_template"`
	Status      string                `json:"status"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type KYCTemplateResponse struct {
	QuestionID string   `json:"question_id"`
	AnswerType string   `json:"answer_type"`
	DependsOn  []string `json:"depends_on"`
	Mandatory  bool     `json:"mandatory"`
}

type UpdateCustomerRequest struct {
	KYCVersion int `json:"kyc_version"`
	KYC        []struct {
		QuestionID string `json:"question_id"`
		Answer     any    `json:"answer"`
	} `json:"kyc"`
}

type UpdateCustomerRequestV2 struct {
	KYCVersion int      `json:"kyc_version"`
	KYC        []string `json:"kyc"`
}

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
