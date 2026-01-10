package dal

import (
	"context"
	"errors"
	"mianshiba/domain/user/internal/dal/model"
	"mianshiba/domain/user/internal/dal/query"
	"time"

	"gorm.io/gorm"
)

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		query: query.Use(db),
	}
}

type UserDAO struct {
	query *query.Query
}

func (dao *UserDAO) CreateUser(ctx context.Context, user *model.User) error {
	return dao.query.User.WithContext(ctx).Create(user)
}

func (dao *UserDAO) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	return dao.query.User.WithContext(ctx).Where(
		dao.query.User.ID.Eq(userID),
	).First()
}

func (dao *UserDAO) GetUserByEmail(ctx context.Context, email string) (*model.User, bool, error) {
	user, err := dao.query.User.WithContext(ctx).Where(dao.query.User.Email.Eq(email)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	return user, true, err
}

func (dao *UserDAO) GetUserByUserName(ctx context.Context, userName string) (*model.User, bool, error) {
	user, err := dao.query.User.WithContext(ctx).Where(dao.query.User.Username.Eq(userName)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	return user, true, err
}

func (dao *UserDAO) UpdateProfile(ctx context.Context, userID int64, updates map[string]interface{}) error {
	if _, ok := updates["updated_at"]; !ok {
		updates["updated_at"] = time.Now()
	}

	_, err := dao.query.User.WithContext(ctx).Where(
		dao.query.User.ID.Eq(userID),
	).Updates(updates)
	return err
}

func (dao *UserDAO) UpdatePassword(ctx context.Context, email, password string) error {
	_, err := dao.query.User.WithContext(ctx).Where(
		dao.query.User.Email.Eq(email),
	).Updates(map[string]interface{}{
		"password":   password,
		"updated_at": time.Now().UnixMilli(),
	})
	return err
}

func (dao *UserDAO) CheckEmailExist(ctx context.Context, email string) (bool, error) {
	_, exist, err := dao.GetUserByEmail(ctx, email)
	if !exist {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
