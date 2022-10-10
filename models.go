package main

import (
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
	CustomerID  uuid.UUID             `json:"customer_id"`
	KYCVersion  int                   `json:"kyc_version"`
	KYCTemplate []KYCTemplateResponse `json:"kyc_template"`
	Status      string                `json:"status"`
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

type UpdateCustomerResponse struct {
	CustomerID uuid.UUID `json:"customer_id"`
	KYCVersion int       `json:"kyc_version"`
	KYC        []struct {
		QuestionID string `json:"question_id"`
		Answer     any    `json:"answer"`
	} `json:"kyc"`
	Status string `json:"status"`
}

type GetCustomerResponse struct {
	CustomerID uuid.UUID `json:"customer_id"`
	KYCVersion int       `json:"kyc_version"`
	KYC        []struct {
		QuestionID string `json:"question_id"`
		Answer     any    `json:"answer"`
	} `json:"kyc"`
	Status string `json:"status"`
}

func newCreateCustomerResponse(customer Customer, kycTemplate []KYCTemplateEntry) CreateCustomerResponse {
	res := CreateCustomerResponse{
		CustomerID: customer.CustomerID,
		KYCVersion: customer.KYCVersion,
		Status:     customer.Status,
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
		CustomerID: customer.CustomerID,
		KYCVersion: customer.KYCVersion,
		Status:     "DRAFT",
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
