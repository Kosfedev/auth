package model

// NewUserData is ...
type NewUserData struct {
	Name     string `redis:"name"`
	Email    string `redis:"email"`
	Role     uint8  `redis:"role"`
	Password string `redis:"password"`
	// TODO: убрать с репо слоя
	PasswordConfirm string `redis:"password_confirm"`
}

// UserData is ...
type UserData struct {
	ID        int64  `redis:"id"`
	Name      string `redis:"name"`
	Email     string `redis:"email"`
	Role      uint8  `redis:"role"`
	CreatedAt int64  `redis:"created_at"`
	UpdatedAt *int64 `redis:"updated_at"`
}

// UpdatedUserData is ...
type UpdatedUserData struct {
	Name  *string `redis:"name"`
	Email *string `redis:"email"`
	Role  *uint8  `redis:"role"`
}
