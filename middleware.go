package tidio

import (
	"context"
	"net/http"
)

type authMid struct {
	service *Service
}

func (m *authMid) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Authorization")
		account, ok := m.service.AccountByKey(key)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "account", account))
		next.ServeHTTP(w, r)
	})
}
