package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/util"
	"net/http"
)

type Handler struct {
	Svc Service
}

func (h Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var u *User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		util.JsonErrorResponse(w, http.StatusBadRequest, fmt.Errorf("parse body: %v", err))
		return
	}

	if err := h.Svc.AddUser(u); err != nil {
		if errors.Is(err, ErrUsernameExists) {
			util.JsonErrorResponse(w, http.StatusUnprocessableEntity, err)
			return
		}

		util.JsonErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	util.JsonResponse(w, http.StatusCreated, u)
}
