package main

import (
	"github.com/google/uuid"
)

var bizSpec = KYCSpecification{
	KYCVersion: 1,
	KYCTemplate: []struct {
		QuestionID string
		AnswerType string
		DependsOn  []string
	}{
		{QuestionID: "external_id_question_1", AnswerType: "text"},                                                           // string type
		{QuestionID: "external_id_question_2", AnswerType: "boolean"},                                                        // true or false
		{QuestionID: "external_id_question_3", AnswerType: "multi_selection", DependsOn: []string{"external_id_question_1"}}, // [a, b, c]
	},
}

type KYCSpecification struct {
	KYCVersion  int
	KYCTemplate []struct {
		QuestionID string
		AnswerType string
		DependsOn  []string
	}
}

type Customer struct {
	CustomerID uuid.UUID
	Answers    []Answer
}

type Answer struct {
	QuestionID string
	Answer     any
}
