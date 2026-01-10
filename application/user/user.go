package user

import (
	"context"
	"mianshiba/api/model/user"
	userAPI "mianshiba/api/model/user"
	"mianshiba/application/base/ctxutil"
	"mianshiba/domain/user/entity"
	"mianshiba/domain/user/service"
	userService "mianshiba/domain/user/service"
	"mianshiba/pkg/errorx"
	"mianshiba/types/errno"
	"net/mail"
)

var UserApplicationSVC = &UserApplicationService{}

type UserApplicationService struct {
	UserDomainSVC      userService.User
	UserModelDomainSVC userService.UserModel
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (u *UserApplicationService) Register(ctx context.Context, req *userAPI.RegisterRequest) (resp *userAPI.LoginResponse, jwtToken string, err error) {
	if !isValidEmail(req.GetEmail()) {
		return nil, "", errorx.New(errno.ErrUserInvalidParamCode, errorx.KV("msg", "Invalid email"))
	}

	userInfo, err := u.UserDomainSVC.Create(ctx, &service.CreateUserRequest{
		UserName: req.GetUsername(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, "", err
	}

	userInfo, jwtToken, err = u.UserDomainSVC.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, "", err
	}

	return &user.LoginResponse{
		Data: userDo2UserTo(userInfo),
		Code: 0,
	}, jwtToken, nil
}

func (u *UserApplicationService) Login(ctx context.Context, req *userAPI.LoginRequest) (resp *userAPI.LoginResponse, jwtToken string, err error) {
	if !isValidEmail(req.GetEmail()) {
		return nil, "", errorx.New(errno.ErrUserInfoInvalidateCode)
	}
	userInfo, jwtToken, err := u.UserDomainSVC.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, "", err
	}

	return &user.LoginResponse{
		Data: userDo2UserTo(userInfo),
		Code: 0,
	}, jwtToken, nil
}

func (u *UserApplicationService) Logout(ctx context.Context, req *userAPI.EmptyRequest) (resp *userAPI.BaseResponse, err error) {
	userID := ctxutil.GetUIDFromCtx(ctx)
	if userID == nil {
		return nil, errorx.New(errno.ErrUserInfoInvalidateCode)
	}

	err = u.UserDomainSVC.Logout(ctx, *userID)
	if err != nil {
		return nil, err
	}

	return &userAPI.BaseResponse{
		Code: 0,
	}, nil
}

// ForgotPassword handles the password reset request
func (u *UserApplicationService) ForgotPassword(ctx context.Context, req *userAPI.ForgotPasswordRequest) (resp *userAPI.EmptyResponse, err error) {
	if !isValidEmail(req.GetEmail()) {
		return nil, errorx.New(errno.ErrUserInvalidParamCode, errorx.KV("msg", "Invalid email"))
	}

	err = u.UserDomainSVC.ForgotPassword(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	// Return an empty response as per the API definition
	// We don't want to reveal if the email exists for security reasons
	return &userAPI.EmptyResponse{}, nil
}

// ResetPassword handles the password reset with token
func (u *UserApplicationService) ResetPassword(ctx context.Context, req *userAPI.ResetPasswordRequest) (resp *userAPI.BaseResponse, err error) {
	if req.GetToken() == "" {
		return nil, errorx.New(errno.ErrUserInvalidParamCode, errorx.KV("msg", "Reset token is required"))
	}

	if len(req.GetPassword()) < 6 {
		return nil, errorx.New(errno.ErrUserInvalidParamCode, errorx.KV("msg", "Password must be at least 6 characters"))
	}

	if req.GetPassword() != req.GetConfirmPassword() {
		return nil, errorx.New(errno.ErrUserInvalidParamCode, errorx.KV("msg", "Passwords do not match"))
	}

	err = u.UserDomainSVC.ResetPassword(ctx, req.GetToken(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &userAPI.BaseResponse{
		Code: 0,
	}, nil
}

func userDo2UserTo(userDo *entity.User) *user.UserProfile {
	return &user.UserProfile{
		ID:        userDo.ID,
		Username:  userDo.Username,
		Email:     userDo.Email,
		Role:      userDo.Role,
		Avatar:    &userDo.Avatar,
		CreatedAt: &userDo.CreatedAt,
		UpdatedAt: &userDo.UpdatedAt,
	}
}
