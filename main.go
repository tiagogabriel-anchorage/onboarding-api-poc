package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)
	configureEndpoints(r)
	log.Fatal(http.ListenAndServe(":3000", r))
}

func configureEndpoints(r *mux.Router) {
	// r.HandleFunc("/", welcome)
	handleRoute(r, "/", welcome, "GET")
	handleRoute(r, "/v1/customers", postCustomers, "POST")
	handleRoute(r, "/v1/customers/{id}", putCustomer, "PUT")

	r.HandleFunc("/v2/customers/{id}", putCustomerV2).Methods("PUT")
	r.HandleFunc("/v1/customers/{id}", getCustomer).Methods("GET")
	r.HandleFunc("/v1/customers/{id}/submit", postCustomerSubmission).Methods("POST")
}
