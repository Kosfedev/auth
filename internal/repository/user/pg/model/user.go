package model

import "time"

// NewUserData is ...
type NewUserData struct {
	Name     string
	Email    string
	Role     uint8
	Password string
}

// UserData is ...
type UserData struct {
	ID        int64
	Name      string
	Email     string
	Role      uint8
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// UpdatedUserData is ...
type UpdatedUserData struct {
	Name  *string
	Email *string
	Role  *uint8
}
