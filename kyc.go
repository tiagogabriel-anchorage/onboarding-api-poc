package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// bizSpec is the KYC specification valid for a specific customer (business or consumer)
// and for a specific market/entity: US, SG, PT, etc.
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

// KYCSpecification is the collection of rules to be applied to a set of customers.
type KYCSpecification struct {
	KYCVersion  int
	KYCTemplate []KYCTemplateEntry
}

// KYCTemplateEntry is a single KYC rule node that constitutes a specification.
type KYCTemplateEntry struct {
	QuestionID string
	AnswerType string
	Mandatory  bool
	DependsOn  []string
}

// getKYCSpec finds a specification valid for a combination of kind of customer and market it belongs to.
func getKYCSpec(kind, entity string) (KYCSpecification, error) {
	if !strings.EqualFold(kind, "business") {
		return KYCSpecification{}, fmt.Errorf("'%s' not supported as customer kind", kind)
	}

	if !strings.EqualFold(entity, "anchorage hold") {
		return KYCSpecification{}, fmt.Errorf("'%s' not supported as entity", entity)
	}

	return bizSpec, nil
}

// db is a memory persistence layer that saves customer's KYC information.
var db = make(map[uuid.UUID]Customer)

// Customer, from a KYC standpoint, is defined by a set of answers that match a specification.
type Customer struct {
	ID         uuid.UUID
	KYCVersion int
	Answers    []Answer
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Answers is the response node of the KYC form, which is defined by the specification.
type Answer struct {
	QuestionID string
	Answer     any
}

// SubmitAnswers provides the customer with new KYC information. Information previously
// set can be overridden, but never deleted, as new information gets added.
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

// saveKYC saves the customer KYC information.
func saveKYC(customer Customer) {
	db[customer.ID] = customer
}

// customerByID finds a customer by its id.
func customerByID(id uuid.UUID) (Customer, bool) {
	c, exists := db[id]
	return c, exists
}
