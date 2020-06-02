package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/gregoryv/stamp"
	"github.com/preferit/tidio"
)

func NewRouter(apikeys map[string]string, store *tidio.Store) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api", serveAPIRoot())

	auth := &authMid{keys: apikeys}
	r.Handle(
		"/api/timesheets/{user}/", auth.Middleware(writeTimesheets(store)),
	).Methods("POST")

	r.Handle(
		"/api/timesheets/", auth.Middleware(readTimesheets()),
	).Methods("GET")
	return r
}

type authMid struct {
	keys map[string]string
}

func (m *authMid) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		if key == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		account, found := m.keys[key]
		if !found {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "account", account))
		next.ServeHTTP(w, r)
	})
}

func readTimesheets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func writeTimesheets(store *tidio.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account, _ := r.Context().Value("account").(string)
		body, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		vars := mux.Vars(r)
		// only allow account to write it's own timesheet
		if vars["user"] != account {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		filename := path.Join(account, "somefile.txt")
		if err := store.WriteFile(filename, body, 0644); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
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
