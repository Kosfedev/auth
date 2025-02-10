package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func parseID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func httpErrorJSON(w http.ResponseWriter, data interface{}, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode create user validation errors", http.StatusInternalServerError)
		return
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}
