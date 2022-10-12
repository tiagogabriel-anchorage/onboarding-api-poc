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

// configureEndpoints sets the endpoints the web service has.
func configureEndpoints(r *mux.Router) {
	handleRoute(r, "/", Welcome, "GET")
	handleRoute(r, "/v1/customers", PostCustomers, "POST")
	handleRoute(r, "/v1/customers/{id}", PutCustomer, "PUT")
	handleRoute(r, "/v2/customers/{id}", PutCustomerV2, "PUT")
	handleRoute(r, "/v1/customers/{id}", GetCustomer, "GET")
	handleRoute(r, "/v1/customers/{id}/submit", PostCustomerSubmission, "POST")
}
