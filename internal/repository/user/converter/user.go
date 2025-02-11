package converter

import (
	"github.com/Kosfedev/auth/internal/model"
	modelRepo "github.com/Kosfedev/auth/internal/repository/user/model"
)

// UserDataFromRepo is...
func UserDataFromRepo(userData *modelRepo.UserData) *model.UserData {
	return &model.UserData{
		ID:        userData.ID,
		Name:      userData.Name,
		Email:     userData.Email,
		Role:      userData.Role,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}
}
