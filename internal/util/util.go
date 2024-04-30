package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var ErrPasswordMismatch = errors.New("passwords don't match")

func HashString(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash string: %v", err)
	}

	return string(hash), nil
}

func CompareHashAndPassword(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return ErrPasswordMismatch
	}

	return nil
}

func JsonErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	log.Printf("error response: status %d %v\n", statusCode, err)

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")

	var msg string
	if statusCode == http.StatusInternalServerError {
		msg = http.StatusText(http.StatusInternalServerError)
	} else {
		msg = err.Error()
	}

	_, _ = w.Write([]byte(msg))
}

func JsonResponse(w http.ResponseWriter, statusCode int, data any) {
	payload, err := json.Marshal(data)
	if err != nil {
		JsonErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("marshal response: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(payload)
}
