package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	newUser := &NewUserData{}
	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		http.Error(w, "Failed to decode new user data", http.StatusBadRequest)
		return
	}

	if err := validateStruct(newUser); err != nil {
		res := &ResponseValidationError{*err}
		httpErrorJSON(w, res, http.StatusBadRequest)
		return
	}

	id, errID := createUser(r.Context(), newUser)
	if errID != nil {
		http.Error(w, "Failed to create new user", http.StatusInternalServerError)
		return
	}

	res := ResponseUserID{ID: id}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode new user id", http.StatusInternalServerError)
		return
	}
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := parseID(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := getUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode new user", http.StatusInternalServerError)
		return
	}
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	updatedUser := &UpdateUserData{}
	userIDStr := chi.URLParam(r, "id")
	userID, err := parseID(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(updatedUser); err != nil {
		http.Error(w, "Failed to decode new user data", http.StatusBadRequest)
		return
	}

	if err := validateStruct(updatedUser); err != nil {
		res := &ResponseValidationError{*err}
		httpErrorJSON(w, res, http.StatusBadRequest)
		return
	}

	res, err := updateUser(r.Context(), updatedUser, userID)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to encode updated user id", http.StatusInternalServerError)
		return
	}
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := parseID(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = deleteUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
