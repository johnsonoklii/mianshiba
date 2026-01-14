package dal

import (
	"context"
	"mianshiba/domain/interview/dal/model"
	"mianshiba/domain/interview/dal/query"

	"gorm.io/gorm"
)

const (
	StatusParsing      = 2 // 解析中
	StatusParseSuccess = 3 // 解析成功
	StatusParseFailed  = 4 // 解析失败
)

func NewResumeDAO(db *gorm.DB) *ResumeDAO {
	return &ResumeDAO{
		query: query.Use(db),
	}
}

type ResumeDAO struct {
	query *query.Query
}

func (r *ResumeDAO) Create(ctx context.Context, resume *model.Resume) error {
	return r.query.Resume.WithContext(ctx).Create(resume)
}

func (r *ResumeDAO) UpdateResume(ctx context.Context, id int64, resume *model.Resume) error {
	_, err := r.query.Resume.WithContext(ctx).Where(
		r.query.Resume.ID.Eq(id),
	).Updates(resume)

	return err
}
