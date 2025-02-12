package types

import (
	"time"
)

// ValidationError is ...
type ValidationError struct {
	Field     string      `json:"field"`
	Tag       string      `json:"tag"`
	TagTarget string      `json:"tagTarget"`
	Value     interface{} `json:"value"`
}

// ResponseValidationError is ...
type ResponseValidationError struct {
	Errors []ValidationError `json:"errors"`
}

// ResponseUserID is ...
type ResponseUserID struct {
	ID int64 `json:"id"`
}

// RequestUpdatedUserData is ...
type RequestUpdatedUserData struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
	Role  *uint8  `json:"role"`
}

// ResponseUserData is ...
type ResponseUserData struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      uint8      `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
