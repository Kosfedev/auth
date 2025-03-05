package converter

import (
	modelService "github.com/Kosfedev/auth/internal/model"
	modelRepo "github.com/Kosfedev/auth/internal/repository/user/redis/model"
)

// UserDataFromRepo is...
func UserDataFromRepo(userData *modelRepo.UserData) *modelService.UserData {
	return &modelService.UserData{
		ID:        userData.ID,
		Name:      userData.Name,
		Email:     userData.Email,
		Role:      userData.Role,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}
}

// UserDataToRepo is...
func UserDataToRepo(userData *modelService.UserData) *modelRepo.UserData {
	return &modelRepo.UserData{
		Name:      userData.Name,
		Email:     userData.Email,
		Role:      userData.Role,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}
}
