package service

import (
	"context"
	"fmt"
	"mianshiba/domain/interview/dal"
	"mianshiba/domain/interview/dal/model"
	"mianshiba/domain/interview/entity"
	"mianshiba/domain/interview/repository"
	"mianshiba/infra/contract/idgen"
	"mianshiba/infra/contract/storage"
	"strings"
)

type ResumeComponents struct {
	OSSClient  storage.Storage
	IDGen      idgen.IDGenerator
	ResumeRepo repository.ResumeRepository
}

func NewResumeDomain(ctx context.Context, c *ResumeComponents) Resume {
	return &resumeImpl{
		ResumeComponents: c,
	}
}

type resumeImpl struct {
	*ResumeComponents
}

func (r *resumeImpl) GetUploadURL(ctx context.Context, userID int64, fileName string, fileType string) (fileID int64, fileKey string, url string, err error) {
	// 生成唯一ID作为文件标识
	fileID, err = r.IDGen.GenID(ctx)
	if err != nil {
		return 0, "", "", err
	}

	// 提取文件扩展名
	ext := ""
	if idx := strings.LastIndex(fileName, "."); idx != -1 {
		ext = fileName[idx:]
	}

	// 构建文件路径（userID/文件标识.扩展名的形式）
	fileKey = fmt.Sprintf("resume/%d/%d%s", userID, fileID, ext)

	// 生成上传URL
	url, err = r.OSSClient.GetUploadUrl(ctx, fileKey)
	if err != nil {
		return 0, "", "", err
	}

	return fileID, fileKey, url, nil
}

func (r *resumeImpl) Create(ctx context.Context, req *ResumeCreateRequest) (resume *entity.Resume, err error) {
	// TODO: 检查文件是否已存在

	// 创建简历实体
	newResume := &model.Resume{
		ID:       req.ID,
		UserID:   req.UserID,
		FileKey:  req.FileKey,
		Filename: req.Filename,
		Filetype: req.Filetype,
		Filesize: req.Filesize,
		Status:   dal.StatusParsing,
	}

	// 保存到数据库
	err = r.ResumeRepo.Create(ctx, newResume)
	if err != nil {
		return nil, err
	}

	return userPo2Do(newResume), nil
}

func userPo2Do(model *model.Resume) *entity.Resume {
	return &entity.Resume{
		ID:          model.ID,
		UserID:      model.UserID,
		FileKey:     model.FileKey,
		Filename:    model.Filename,
		Filetype:    model.Filetype,
		Filesize:    model.Filesize,
		Status:      model.Status,
		ParseStatus: model.ParseStatus,
		ParseError:  model.ParseError,
		UploadAt:    model.UploadAt.UnixMilli(),
	}
}
