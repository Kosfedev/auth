package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	userImplementation "github.com/Kosfedev/auth/internal/api/user"
	"github.com/Kosfedev/auth/internal/model"
	"github.com/Kosfedev/auth/pkg/user_v1/http/types"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// TODO: Relocate
var validate = validator.New()

func validateStruct(data interface{}) *[]types.ValidationError {
	var errors []types.ValidationError

	err := validate.Struct(data)
	if err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, err := range errs {
				errors = append(errors, types.ValidationError{
					Field:     err.Field(),
					Tag:       err.Tag(),
					TagTarget: err.Param(),
					Value:     err.Value(),
				})
			}
		}

		return &errors
	}

	return nil
}

// CreateUserHandler is...
func CreateUserHandler(w http.ResponseWriter, r *http.Request, userServiceImpl userImplementation.Implementation) {
	user := &model.NewUserData{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		http.Error(w, "Failed to decode new userImplementation data", http.StatusBadRequest)
		return
	}

	if err := validateStruct(user); err != nil {
		res := &types.ResponseValidationError{Errors: *err}
		httpErrorJSON(w, res, http.StatusBadRequest)
		return
	}

	resWithID, errID := userServiceImpl.Create(r.Context(), user)
	if errID != nil {
		http.Error(w, "Failed to create new userImplementation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resWithID); err != nil {
		http.Error(w, "Failed to encode new userImplementation id", http.StatusInternalServerError)
		return
	}
}

// GetUserHandler is...
func GetUserHandler(w http.ResponseWriter, r *http.Request, userServiceImpl userImplementation.Implementation) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid userImplementation ID", http.StatusBadRequest)
		return
	}

	user, err := userServiceImpl.Get(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Failed to encode new userImplementation", http.StatusInternalServerError)
		return
	}
}

// PutUserHandler is...
func PutUserHandler(w http.ResponseWriter, r *http.Request, userServiceImpl userImplementation.Implementation) {
	updatedUser := &types.RequestUpdatedUserData{}
	// TODO: const?
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid userImplementation ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(updatedUser); err != nil {
		http.Error(w, "Failed to decode new userImplementation data", http.StatusBadRequest)
		return
	}

	if err := validateStruct(updatedUser); err != nil {
		res := &types.ResponseValidationError{Errors: *err}
		httpErrorJSON(w, res, http.StatusBadRequest)
		return
	}

	resUpdatedUser, err := userServiceImpl.Patch(r.Context(), updatedUser, userID)
	if err != nil {
		http.Error(w, "Failed to update userImplementation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resUpdatedUser); err != nil {
		http.Error(w, "Failed to encode updated userImplementation id", http.StatusInternalServerError)
		return
	}
}

// DeleteUserHandler is...
func DeleteUserHandler(w http.ResponseWriter, r *http.Request, userServiceImpl userImplementation.Implementation) {
	// TODO: const?
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid userImplementation ID", http.StatusBadRequest)
		return
	}

	err = userServiceImpl.Delete(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to delete userImplementation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func httpErrorJSON(w http.ResponseWriter, data interface{}, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode create userImplementation validation errors", http.StatusInternalServerError)
		return
	}
}
