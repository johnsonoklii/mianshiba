package service

import (
	"context"
	"mianshiba/domain/interview/entity"
)

type ResumeCreateRequest struct {
	ID       int64
	UserID   int64
	FileKey  string
	Filename string
	Filetype string
	Filesize int64
}

type Resume interface {
	GetUploadURL(ctx context.Context, userID int64, fileName string, fileType string) (fileID int64, fileKey string, url string, err error)
	Create(ctx context.Context, req *ResumeCreateRequest) (resume *entity.Resume, err error)
}
