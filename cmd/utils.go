package main

import (
	"encoding/json"
	"net/http"
)

func httpErrorJSON(w http.ResponseWriter, data interface{}, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode create user validation errors", http.StatusInternalServerError)
		return
	}
}
