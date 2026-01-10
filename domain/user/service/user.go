package service

import (
	"context"
	"mianshiba/domain/user/entity"
)

type CreateUserRequest struct {
	UserName string
	Email    string
	Password string
	Avatar   string
}

type UpdateProfileRequest struct {
}

type User interface {
	Create(ctx context.Context, req *CreateUserRequest) (user *entity.User, err error)
	Login(ctx context.Context, email, password string) (user *entity.User, jwtToken string, err error)
	Logout(ctx context.Context, userID int64) (err error)
	ForgotPassword(ctx context.Context, email string) (err error)
	ResetPassword(ctx context.Context, token, password string) (err error)
	GetUser(ctx context.Context, userID int64) (user *entity.User, err error)
	GetUserList(ctx context.Context) (userList []*entity.User, err error)
	UpdateProfile(ctx context.Context, req *UpdateProfileRequest) (err error)
	GetJwtToken(ctx context.Context, userID int64) (jwtToken string, err error)
}
