package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gregoryv/stamp"
	"github.com/preferit/tidio"
)

func NewRouter(service *tidio.Service) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api", serveAPIRoot())

	auth := (&authMid{service: service}).Middleware
	r.Handle(
		"/api/timesheets/{user}/{filename}", auth(writeTimesheets()),
	).Methods("POST")
	r.Handle(
		"/api/timesheets/{user}/{filename}", auth(readTimesheets()),
	).Methods("GET")
	r.Handle(
		"/api/timesheets/{user}/", auth(listTimesheets()),
	).Methods("GET")
	return r
}

func writeTimesheets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value("role").(*tidio.Role)
		vars := mux.Vars(r)
		filename := vars["filename"]
		user := vars["user"]
		s := &tidio.Timesheet{
			Filename: filename,
			Owner:    user,
			Content:  r.Body,
		}
		if err := role.CreateTimesheet(s); err != nil {
			w.WriteHeader(statusOf(err))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func readTimesheets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value("role").(*tidio.Role)
		vars := mux.Vars(r)
		filename := vars["filename"]
		user := vars["user"]
		var buf bytes.Buffer
		if err := role.ReadTimesheet(&buf, filename, user); err != nil {
			w.WriteHeader(statusOf(err))
			return
		}
		fmt.Println(buf.String())
		buf.WriteTo(w)
	}
}

func listTimesheets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			role, _ = r.Context().Value("role").(*tidio.Role)
			vars    = mux.Vars(r)
			user    = vars["user"]
		)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"timesheets": role.ListTimesheet(user),
		})
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

func statusOf(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case err == tidio.ErrForbidden:
		return http.StatusForbidden
	default:
		return http.StatusBadRequest
	}
}
