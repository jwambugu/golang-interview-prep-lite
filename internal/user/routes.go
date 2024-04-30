package user

import (
	"errors"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/util"
	"net/http"
	"strings"
)

func (h Handler) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			util.JsonErrorResponse(w, http.StatusUnauthorized, errors.New("auth header is required"))
			return
		}

		headerValues := strings.Split(authHeader, " ")
		if len(headerValues) != 2 {
			util.JsonErrorResponse(w, http.StatusUnauthorized, errors.New("invalid auth header provided"))
			return
		}

		claims, err := h.jwtManager.Validate(headerValues[1])
		if err != nil {
			util.JsonErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		if _, err = h.repo.FindByID(r.Context(), claims.ID); err != nil {
			if errors.Is(err, ErrRecordNotFound) {
				util.JsonErrorResponse(w, http.StatusNotFound, err)
				return
			}

			util.JsonErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/login", h.Authenticate)
	mux.Handle("/user", h.requireAuth(http.HandlerFunc(h.AddUser)))
}
