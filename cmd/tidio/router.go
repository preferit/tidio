package main

import (
	"encoding/json"
	"io"
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
		vars := mux.Vars(r)
		filename := vars["filename"]
		s := &tidio.Timesheet{
			Path:       filename,
			ReadCloser: r.Body,
		}
		role, _ := r.Context().Value("role").(*tidio.Role)
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
		sheet := tidio.Timesheet{
			Path: vars["filename"],
		}
		if err := role.OpenTimesheet(&sheet); err != nil {
			w.WriteHeader(statusOf(err))
			return
		}
		io.Copy(w, sheet)
		sheet.Close()
	}
}

func listTimesheets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value("role").(*tidio.Role)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"timesheets": role.ListTimesheet(),
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
