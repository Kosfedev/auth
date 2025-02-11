package user

import (
	"context"
	"log"

	"github.com/Kosfedev/auth/internal/converter"
	modelHTTP "github.com/Kosfedev/auth/pkg/user_v1"
)

// Get is...
func (i *Implementation) Get(ctx context.Context, id int64) (*modelHTTP.ResponseUserData, error) {
	user, err := i.userService.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, name: %s, email: %s, role: %v, created_at: %v, updated_at: %v\n", user.ID, user.Name, user.Email, user.Role, user.CreatedAt, user.UpdatedAt)

	return converter.UserDataFromServiceToHTTP(user), nil
}
