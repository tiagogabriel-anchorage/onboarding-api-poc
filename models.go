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

type NewCustomerRequest struct {
	CustomerKind string `json:"kind"`
	Entity       string `json:"entity"`
}

type NewCustomerResponse struct {
	CustomerID  uuid.UUID             `json:"customer_id"`
	KYCVersion  int                   `json:"kyc_version"`
	KYCTemplate []KYCTemplateResponse `json:"kyc_template"`
}

type KYCTemplateResponse struct {
	QuestionID string   `json:"question_id"`
	AnswerType string   `json:"answer_type"`
	DependsOn  []string `json:"depends_on"`
	Mandatory  bool     `json:"mandatory"`
}

func newCustomerResponse(customer Customer, kycTemplate []KYCTemplateEntry) NewCustomerResponse {
	res := NewCustomerResponse{
		CustomerID: customer.CustomerID,
		KYCVersion: customer.KYCVersion,
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
