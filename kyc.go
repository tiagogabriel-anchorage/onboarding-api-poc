package main

import (
	"fmt"
	"strings"
	"time"

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
	ID         uuid.UUID
	KYCVersion int
	Answers    []Answer
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Answer struct {
	QuestionID string
	Answer     any
}

func (c *Customer) SubmitAnswers(newAnswers []Answer) {
	currentAnswers := map[string]Answer{}
	for _, answer := range c.Answers {
		currentAnswers[answer.QuestionID] = answer
	}

	// submit new answers (override if exists, create if new)
	for _, answer := range newAnswers {
		currentAnswers[answer.QuestionID] = answer
	}

	c.Answers = nil // forget about previous answers
	// set new answers
	for _, answer := range currentAnswers {
		c.Answers = append(c.Answers, answer)
	}
}

func saveKYC(customer Customer) {
	db[customer.ID] = customer
}

func customerById(id uuid.UUID) (Customer, bool) {
	c, exists := db[id]
	return c, exists
}
