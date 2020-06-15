package tidio

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(service *Service) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api", serveAPIRoot(service))

	auth := service.FindAccount
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
		resource := Resource{
			Path:       vars["filename"],
			ReadCloser: r.Body,
		}
		account, _ := r.Context().Value("account").(*Account)
		if err := account.WriteResource(&resource); err != nil {
			w.WriteHeader(statusOf(err))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func readTimesheets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account, _ := r.Context().Value("account").(*Account)
		vars := mux.Vars(r)
		sheet := Timesheet{
			Path: vars["filename"],
		}
		if err := account.OpenTimesheet(&sheet); err != nil {
			w.WriteHeader(statusOf(err))
			return
		}
		io.Copy(w, sheet)
		sheet.Close()
	}
}

func listTimesheets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account, _ := r.Context().Value("account").(*Account)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"timesheets": account.ListTimesheet(),
		})
	}
}

func serveAPIRoot(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resource := Resource{Path: ""}
		service.FindResource(&resource)
		io.Copy(w, resource)
		resource.Close()
	}
}

func statusOf(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case err == ErrForbidden:
		return http.StatusForbidden
	default:
		return http.StatusBadRequest
	}
}