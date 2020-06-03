package main

import (
	"context"
	"encoding/json"
	"net/http"

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
		vars := mux.Vars(r)
		filename := vars["filename"]
		user := vars["user"]
		if err := role.CreateTimesheet(filename, user, r.Body); err != nil {
			w.WriteHeader(statusOf(err))
			return
		}
		w.WriteHeader(http.StatusNoContent)
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
