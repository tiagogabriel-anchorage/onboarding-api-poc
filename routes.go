package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJson(w http.ResponseWriter, statusCode int, body any) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Fatalf("Could not parse the response. Err: %s", err)
	}
	w.WriteHeader(statusCode)
	w.Write(jsonBody)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusOK, map[string]string{
		"message": "Welcome to Onboarding API (PoC)",
	})
}

func postCustomers(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, http.StatusCreated, map[string]string{
		"spec": "dummy",
	})
}
