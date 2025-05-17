package model

import "time"

// NewUserData is ...
type NewUserData struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Role            uint8  `json:"role" validate:"required"`
	Password        string `json:"password" validate:"required,min=5"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
}

// UserData is ...
type UserData struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      uint8      `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// UpdatedUserData is ...
type UpdatedUserData struct {
	Name *string `json:"name" validate:"required_without_all=Email Role,min=1"`
	// TODO: email error even if not required field
	Email *string `json:"email" validate:"required_without_all=Name Role,email"`
	Role  *uint8  `json:"role" validate:"required_without_all=Name Email"`
}
