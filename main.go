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
	handleRoute(r, "/", welcome, "GET")
	handleRoute(r, "/v1/customers", postCustomers, "POST")
	handleRoute(r, "/v1/customers/{id}", putCustomer, "PUT")
	handleRoute(r, "/v2/customers/{id}", putCustomerV2, "PUT")
	handleRoute(r, "/v1/customers/{id}", getCustomer, "GET")

	r.HandleFunc("/v1/customers/{id}/submit", postCustomerSubmission).Methods("POST")
}
