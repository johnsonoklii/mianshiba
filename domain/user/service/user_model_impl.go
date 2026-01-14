package service

import (
	"context"
	"fmt"
	"mianshiba/domain/user/dal/model"
	"mianshiba/domain/user/entity"
	"mianshiba/domain/user/repository"
	"mianshiba/infra/contract/cache"
	"mianshiba/infra/contract/idgen"
	"mianshiba/pkg/encrypt"
	"time"
)

type UserModelComponents struct {
	CacheCli      cache.Cmdable
	IDGen         idgen.IDGenerator
	UserModelRepo repository.UserModelRepository
}

func NewUserModelDomain(ctx context.Context, c *UserModelComponents) UserModel {
	return &userModelImpl{
		UserModelComponents: c,
	}
}

type userModelImpl struct {
	*UserModelComponents
}

func (u *userModelImpl) Create(ctx context.Context, userID int64, req *CreateUserModelRequest) (userModel *entity.UserModel, err error) {
	apiKey, err := encrypt.EncryptAPIKey(req.APIKey)
	if err != nil {
		return nil, fmt.Errorf("encrypt api key failed: %w", err)
	}

	secretHint := ""
	if len(req.APIKey) > 4 {
		secretHint = "***" + req.APIKey[len(req.APIKey)-4:]
	} else {
		secretHint = "***" + req.APIKey
	}

	modelID, err := u.IDGen.GenID(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate id error: %w", err)
	}

	newModel := &model.UserModel{
		ID:              modelID,
		UserID:          userID,
		Name:            req.Name,
		ModelKey:        req.ModelKey,
		Protocol:        req.Protocol,
		BaseURL:         req.BaseURL,
		APIKeyEncrypted: apiKey,
		SecretHint:      secretHint,
		ConfigJSON:      req.ConfigJSON,
		MetaID:          req.MetaID,
		DefaultParams:   req.DefaultParams,
		Scope:           req.Scope,
		Status:          req.Status,
		IsDefault:       req.IsDefault,
		ProviderName:    req.ProviderName,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Deleted:         false,
	}

	err = u.UserModelRepo.CreateUserModel(ctx, newModel)
	if err != nil {
		return nil, err
	}

	return userModelPo2Do(newModel), nil
}

func userModelPo2Do(model *model.UserModel) *entity.UserModel {
	return &entity.UserModel{
		ID:            model.ID,
		UserID:        model.UserID,
		Name:          model.Name,
		ModelKey:      model.ModelKey,
		Protocol:      model.Protocol,
		BaseURL:       model.BaseURL,
		ConfigJSON:    model.ConfigJSON,
		MetaID:        model.MetaID,
		DefaultParams: model.DefaultParams,
		SecretHint:    model.SecretHint,
		Scope:         model.Scope,
		Status:        model.Status,
		IsDefault:     model.IsDefault,
		ProviderName:  model.ProviderName,
		CreatedAt:     model.CreatedAt.UnixMilli(),
		UpdatedAt:     model.UpdatedAt.UnixMilli(),
	}
}
