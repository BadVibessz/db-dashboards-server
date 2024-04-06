package mapper

import (
	"db-dashboards/internal/domain/entity"
	"db-dashboards/internal/handler/request"
	"db-dashboards/internal/handler/response"
)

func MapUserToUserResponse(user *entity.User) response.GetUserResponse {
	return response.GetUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func MapRegisterRequestToUserEntity(registerReq *request.RegisterRequest) entity.User {
	return entity.User{
		Email:          registerReq.Email,
		HashedPassword: registerReq.Password,
	}
}
