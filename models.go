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
	CustomerType string `json:"type"`
	Entity       string `json:"entity"`
}

type NewCustomerResponse struct {
	CustomerID  uuid.UUID `json:"customer_id"`
	KYCVersion  int       `json:"kyc_version"`
	KYCTemplate []struct {
		QuestionID string   `json:"question_id"`
		AnswerType string   `json:"answer_type"`
		DependsOn  []string `json:"depends_on"`
		Answer     any      `json:"answer"`
	} `json:"kyc_template"`
}
