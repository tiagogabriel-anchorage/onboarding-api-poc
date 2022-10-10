package main

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Specs

var bizSpec = KYCSpecification{
	KYCVersion: 1,
	KYCTemplate: []KYCTemplateEntry{
		{
			QuestionID: "external_id_question_1",
			AnswerType: "text",
			Mandatory:  true,
			DependsOn:  []string{},
		},
		{
			QuestionID: "external_id_question_2",
			AnswerType: "boolean",
			Mandatory:  false,
			DependsOn:  []string{},
		},
		{
			QuestionID: "external_id_question_3",
			AnswerType: "multi_selection",
			Mandatory:  true,
			DependsOn:  []string{"external_id_question_1"},
		},
	},
}

type KYCSpecification struct {
	KYCVersion  int
	KYCTemplate []KYCTemplateEntry
}

type KYCTemplateEntry struct {
	QuestionID string
	AnswerType string
	Mandatory  bool
	DependsOn  []string
}

func getKYCSpec(kind, entity string) (KYCSpecification, error) {
	if !strings.EqualFold(kind, "business") {
		return KYCSpecification{}, fmt.Errorf("'%s' not supported as customer kind", kind)
	}

	if !strings.EqualFold(entity, "anchorage hold") {
		return KYCSpecification{}, fmt.Errorf("'%s' not supported as entity", entity)
	}

	return bizSpec, nil
}

// Customers' answers

var db = make(map[uuid.UUID]Customer)

type Customer struct {
	CustomerID uuid.UUID
	KYCVersion int
	Answers    []Answer
	Status     string
}

type Answer struct {
	QuestionID string
	Answer     any
}

func saveKYC(customer Customer) {
	db[customer.CustomerID] = customer
}

func customerById(id uuid.UUID) Customer {
	if entry, exists := db[id]; exists {
		return entry
	}
	return Customer{}
}
