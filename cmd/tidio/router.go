package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	version  = "0.0"
	revision = "dev"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"revision": revision,
			"version":  version,
		})
	})
	return r
}
