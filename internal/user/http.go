package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/model"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/util"
	"net/http"
)

type Handler struct {
	Svc Service
}

func (h Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req *AuthenticateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.JsonErrorResponse(w, http.StatusBadRequest, fmt.Errorf("parse request: %v", err))
		return
	}

	resp, err := h.Svc.Authenticate(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			util.JsonErrorResponse(w, http.StatusUnauthorized, err)
			return
		}

		util.JsonErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.JsonResponse(w, http.StatusOK, resp)
}

func (h Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var u *model.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		util.JsonErrorResponse(w, http.StatusBadRequest, fmt.Errorf("parse request: %v", err))
		return
	}

	if err := h.Svc.Create(r.Context(), u); err != nil {
		if errors.Is(err, ErrUsernameExists) {
			util.JsonErrorResponse(w, http.StatusUnprocessableEntity, err)
			return
		}

		util.JsonErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.JsonResponse(w, http.StatusCreated, u)
}
