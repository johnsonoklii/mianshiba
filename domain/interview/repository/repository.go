package repository

import (
	"context"
	"mianshiba/domain/interview/dal"
	"mianshiba/domain/interview/dal/model"

	"gorm.io/gorm"
)

func NewResumeRepo(db *gorm.DB) ResumeRepository {
	return dal.NewResumeDAO(db)
}

type ResumeRepository interface {
	Create(ctx context.Context, resume *model.Resume) error
	UpdateResume(ctx context.Context, id int64, resume *model.Resume) error
}
