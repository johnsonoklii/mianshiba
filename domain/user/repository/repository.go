package repository

import (
	"context"
	"mianshiba/domain/user/internal/dal"
	"mianshiba/domain/user/internal/dal/model"

	"gorm.io/gorm"
)

func NewUserRepo(db *gorm.DB) UserRepository {
	return dal.NewUserDAO(db)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, userID int64) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, bool, error)
	GetUserByUserName(ctx context.Context, userName string) (*model.User, bool, error)
	UpdateProfile(ctx context.Context, userID int64, updates map[string]interface{}) error
	UpdatePassword(ctx context.Context, email, password string) error
	CheckEmailExist(ctx context.Context, email string) (bool, error)
}

// UserModel
func NewUserModelRepo(db *gorm.DB) UserModelRepository {
	return dal.NewUserModelDAO(db)
}

type UserModelRepository interface {
	CreateUserModel(ctx context.Context, userModel *model.UserModel) error
	GetUserModelByID(ctx context.Context, userModelID int64) (*model.UserModel, error)
}
