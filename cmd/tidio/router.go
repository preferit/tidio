package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gregoryv/stamp"
)

func NewRouter(apikeys map[string]string) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api", serveAPIRoot())

	rt := r.Path("/api/timesheets/").Subrouter()

	rt.HandleFunc("/", readTimesheets()).Methods("GET")
	rt.HandleFunc("/", writeTimesheets()).Methods("POST")

	auth := &authMid{keys: apikeys}
	rt.Use(auth.Middleware)
	return r
}

type authMid struct {
	keys map[string]string
}

func (m *authMid) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		_, found := m.keys[key]
		if key == "" || !found {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func readTimesheets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func writeTimesheets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func serveAPIRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"revision": stamp.InUse().Revision,
			"version":  stamp.InUse().ChangelogVersion,
			"resources": []string{
				"/api/timesheets/",
			},
		})
	}
}
