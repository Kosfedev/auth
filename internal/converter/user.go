package converter

import (
	"github.com/Kosfedev/auth/internal/model"
	modelHTTP "github.com/Kosfedev/auth/pkg/user_v1"
)

func UpdatedUserDataFromHTTPToService(userData *modelHTTP.RequestUpdatedUserData) *model.UpdatedUserData {
	return &model.UpdatedUserData{
		Name:  userData.Name,
		Email: userData.Email,
		Role:  userData.Role,
	}
}

func UserDataFromServiceToHTTP(userData *model.UserData) *modelHTTP.ResponseUserData {
	return &modelHTTP.ResponseUserData{
		ID:        userData.ID,
		Name:      userData.Name,
		Email:     userData.Email,
		Role:      userData.Role,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}
}
