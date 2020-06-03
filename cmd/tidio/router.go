package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/gregoryv/stamp"
	"github.com/preferit/tidio"
)

func NewRouter(store *tidio.Store, service *tidio.Service) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api", serveAPIRoot())

	auth := &authMid{
		service: service,
	}
	r.Handle(
		"/api/timesheets/{user}/{filename}",
		auth.Middleware(writeTimesheets(store)),
	).Methods("POST")

	r.Handle(
		"/api/timesheets/", auth.Middleware(readTimesheets()),
	).Methods("GET")
	return r
}

type authMid struct {
	service *tidio.Service
}

func (m *authMid) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		role, ok := m.service.IsAuthenticated(key)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "role", role))
		next.ServeHTTP(w, r)
	})
}

func writeTimesheets(store *tidio.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value("role").(*tidio.Role)
		account := role.Account()
		body, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		vars := mux.Vars(r)
		filename := vars["filename"]
		if err := checkTimesheetFilename(filename); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// only allow account to write it's own timesheet
		if vars["user"] != account {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if err := store.WriteFile(filename, body, 0644); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func checkTimesheetFilename(name string) error {
	format := `\d\d\d\d\d\d\.timesheet`
	if ok, _ := regexp.MatchString(format, name); !ok {
		return fmt.Errorf("bad filename: expected format %s", format)
	}
	return nil
}

func readTimesheets() http.HandlerFunc {
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
