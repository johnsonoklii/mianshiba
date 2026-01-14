package dal

import (
	"context"
	"mianshiba/domain/user/dal/model"
	"mianshiba/domain/user/dal/query"

	"gorm.io/gorm"
)

func NewUserModelDAO(db *gorm.DB) *UserModelDAO {
	return &UserModelDAO{
		query: query.Use(db),
	}
}

type UserModelDAO struct {
	query *query.Query
}

func (u *UserModelDAO) CreateUserModel(ctx context.Context, userModel *model.UserModel) error {
	return u.query.UserModel.WithContext(ctx).Create(userModel)
}

func (u *UserModelDAO) GetUserModelByID(ctx context.Context, userModelID int64) (*model.UserModel, error) {
	return u.query.UserModel.WithContext(ctx).Where(u.query.UserModel.ID.Eq(userModelID)).First()
}
