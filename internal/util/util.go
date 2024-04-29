package util

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func HashString(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash string: %v", err)
	}

	return string(hash), nil
}

func CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
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

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(payload)
}
