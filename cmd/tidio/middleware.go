package main

import (
	"context"
	"net/http"

	"github.com/preferit/tidio"
)

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
