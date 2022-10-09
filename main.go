package main

import (
	"fmt"
	"log"
	"net/http"
)

// Hello world
func Welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Onboarding API")
}

func main() {
	http.HandleFunc("/", Welcome)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
