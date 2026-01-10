package service

import (
	"context"
	"mianshiba/domain/user/entity"
)

type CreateUserModelRequest struct {
	Name          string `json:"name"`
	ModelKey      string `json:"model_key"`
	Protocol      string `json:"protocol"`
	BaseURL       string `json:"base_url"`
	APIKey        string `json:"api_key"`
	ProviderName  string `json:"provider_name"`
	MetaID        int64  `json:"meta_id"`
	DefaultParams string `json:"default_params"`
	ConfigJSON    string `json:"config_json"`
	Scope         int32  `json:"scope"`
	Status        int32  `json:"status"`
	IsDefault     int32  `json:"is_default"`
}

type UserModel interface {
	Create(ctx context.Context, userID int64, req *CreateUserModelRequest) (userModel *entity.UserModel, err error)
}
