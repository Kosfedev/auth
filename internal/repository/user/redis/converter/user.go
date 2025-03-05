package converter

import (
	"time"

	modelService "github.com/Kosfedev/auth/internal/model"
	modelRepo "github.com/Kosfedev/auth/internal/repository/user/redis/model"
)

// UserDataFromRepo is...
func UserDataFromRepo(userData *modelRepo.UserData) *modelService.UserData {
	var updatedAt *time.Time
	if userData.UpdatedAt != nil && *userData.UpdatedAt != 0 {
		temp := time.Unix(0, *userData.UpdatedAt)
		updatedAt = &temp
	}

	return &modelService.UserData{
		ID:        userData.ID,
		Name:      userData.Name,
		Email:     userData.Email,
		Role:      userData.Role,
		CreatedAt: time.Unix(0, userData.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

// UserDataToRepo is...
func UserDataToRepo(userData *modelService.UserData) *modelRepo.UserData {
	var updatedAt *int64
	if userData.UpdatedAt != nil {
		temp := (*userData.UpdatedAt).Unix()
		updatedAt = &temp
	}

	return &modelRepo.UserData{
		Name:      userData.Name,
		Email:     userData.Email,
		Role:      userData.Role,
		CreatedAt: userData.CreatedAt.Unix(),
		UpdatedAt: updatedAt,
	}
}
