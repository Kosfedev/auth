package main

import "time"

// Role is ...
type Role uint8

// NewUserData is ...
type NewUserData struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	// TODO: role enum validation
	Role            uint8  `json:"role" validate:"required"`
	Password        string `json:"password" validate:"required,min=5"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
}

// UserData is ...
type UserData struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      Role       `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// UpdateUserData is ...
type UpdateUserData struct {
	Name  *string `json:"name" validate:"required_without_all=Email Role,min=1"`
	Email *string `json:"email" validate:"required_without_all=Name Role,email"`
	// TODO: role enum validation
	Role *Role `json:"role" validate:"required_without_all=Name Email"`
}

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
