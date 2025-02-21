package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Kosfedev/auth/internal/model"
	modelHTTP "github.com/Kosfedev/auth/pkg/user_v1"
	"github.com/go-chi/chi"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	user := &model.NewUserData{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, "Failed to decode new user data", http.StatusBadRequest)
		return
	}

	if err := validateStruct(user); err != nil {
		res := &modelHTTP.ResponseValidationError{Errors: *err}
		httpErrorJSON(w, res, http.StatusBadRequest)
		return
	}

	resWithID, errID := userServiceImpl.Create(r.Context(), user)
	if errID != nil {
		http.Error(w, "Failed to create new user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resWithID); err != nil {
		http.Error(w, "Failed to encode new user id", http.StatusInternalServerError)
		return
	}
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := userServiceImpl.Get(r.Context(), userID)
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
	updatedUser := &modelHTTP.RequestUpdatedUserData{}
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(updatedUser); err != nil {
		http.Error(w, "Failed to decode new user data", http.StatusBadRequest)
		return
	}

	if err := validateStruct(updatedUser); err != nil {
		res := &modelHTTP.ResponseValidationError{Errors: *err}
		httpErrorJSON(w, res, http.StatusBadRequest)
		return
	}

	resUpdatedUser, err := userServiceImpl.Patch(r.Context(), updatedUser, userID)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resUpdatedUser); err != nil {
		http.Error(w, "Failed to encode updated user id", http.StatusInternalServerError)
		return
	}
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = userServiceImpl.Delete(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
