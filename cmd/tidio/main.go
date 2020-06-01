package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "tidio API")
	})
	log.Fatal(http.ListenAndServe(":13001", nil))
}
