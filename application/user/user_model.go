package user

import (
	"context"
	userAPI "mianshiba/api/model/user"
	"mianshiba/application/base/ctxutil"
	"mianshiba/domain/user/entity"
	"mianshiba/domain/user/service"
	"mianshiba/pkg/errorx"
	"mianshiba/types/errno"
)

func (u *UserApplicationService) CreateUserModel(ctx context.Context, req *userAPI.CreateUserModelRequest) (resp *userAPI.CreateUserModelResponse, err error) {
	// Get user ID from context
	userID := ctxutil.GetUIDFromCtx(ctx)
	if userID == nil {
		return nil, errorx.New(errno.ErrUserInfoInvalidateCode, errorx.KV("msg", "User not authenticated"))
	}

	var MetaID int64
	if req.MetaID != nil {
		MetaID = *req.MetaID
	}

	var DefaultParams = "{}"
	if req.DefaultParams != nil {
		DefaultParams = *req.DefaultParams
	}

	var ConfigJSON = "{}"
	if req.ConfigJSON != nil {
		ConfigJSON = *req.ConfigJSON
	}

	var Scope int32
	if req.Scope != nil {
		Scope = *req.Scope
	}

	var Status int32 = 1
	if req.Status != nil {
		Status = *req.Status
	}

	var IsDefault int32
	if req.IsDefault != nil {
		IsDefault = *req.IsDefault
	}

	userModel, err := u.UserModelDomainSVC.Create(ctx, *userID, &service.CreateUserModelRequest{
		Name:          req.Name,
		ModelKey:      req.ModelKey,
		Protocol:      req.Protocol,
		BaseURL:       req.BaseURL,
		APIKey:        req.APIKey,
		ProviderName:  req.ProviderName,
		MetaID:        MetaID,
		DefaultParams: DefaultParams,
		ConfigJSON:    ConfigJSON,
		Scope:         Scope,
		Status:        Status,
		IsDefault:     IsDefault,
	})
	if err != nil {
		return nil, err
	}

	return &userAPI.CreateUserModelResponse{
		Data: userModelDo2To(userModel),
		Code: 0,
	}, nil
}

func userModelDo2To(userModel *entity.UserModel) *userAPI.UserModelDetail {
	return &userAPI.UserModelDetail{
		ID:            userModel.ID,
		UserID:        userModel.UserID,
		Name:          userModel.Name,
		ModelKey:      userModel.ModelKey,
		Protocol:      userModel.Protocol,
		BaseURL:       userModel.BaseURL,
		ConfigJSON:    &userModel.ConfigJSON,
		MetaID:        &userModel.MetaID,
		DefaultParams: &userModel.DefaultParams,
		Scope:         userModel.Scope,
		Status:        userModel.Status,
		SecretHint:    &userModel.SecretHint,
		IsDefault:     userModel.IsDefault,
		ProviderName:  userModel.ProviderName,
		CreatedAt:     userModel.CreatedAt,
		UpdatedAt:     userModel.UpdatedAt,
	}
}
