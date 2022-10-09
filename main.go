package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Onboarding API")
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
