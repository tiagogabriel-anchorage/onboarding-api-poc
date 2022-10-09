package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", welcome)
	http.HandleFunc("/v1/customers", postCustomers)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
