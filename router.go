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
		"/api/timesheets/{user}/{filename}", auth(writeResource()),
	).Methods("POST")
	r.Handle(
		"/api/timesheets/{user}/{filename}", auth(readResource()),
	).Methods("GET")
	r.Handle(
		"/api/timesheets/{user}/", auth(readResources()),
	).Methods("GET")
	return r
}

func writeResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		account, _ := r.Context().Value("account").(*Account)
		if err := account.WriteResource(vars["filename"], r.Body); err != nil {
			w.WriteHeader(statusOf(err))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func readResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account, _ := r.Context().Value("account").(*Account)
		vars := mux.Vars(r)
		resource := Resource{
			Path: vars["filename"],
		}
		if err := account.ReadResource(&resource); err != nil {
			w.WriteHeader(statusOf(err))
			return
		}
		io.Copy(w, resource)
		resource.Close()
	}
}

func readResources() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account, _ := r.Context().Value("account").(*Account)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"timesheets": account.FindResources(),
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
