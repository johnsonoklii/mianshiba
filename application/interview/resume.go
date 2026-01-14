package interview

import (
	"context"
	"encoding/json"
	"time"

	interviewAPI "mianshiba/api/model/interview"
	"mianshiba/application/base/ctxutil"
	"mianshiba/conf"
	"mianshiba/domain/interview/entity"
	"mianshiba/domain/interview/service"
	"mianshiba/pkg/errorx"
	"mianshiba/pkg/logs"
	"mianshiba/types/errno"
)

// ResumeMsg Kafka消息结构
type ResumeMsg struct {
	FileKey  string `json:"file_key"`
	FileID   int64  `json:"file_id"`
	Filename string `json:"filename"`
	Filetype string `json:"filetype"`
	Filesize int64  `json:"filesize"`
	UserID   int64  `json:"user_id"` // 添加用户ID
}

func (i *InterviewApplicationService) GetResumeUploadUrl(ctx context.Context, fileName string, fileType string) (res *interviewAPI.ResumeUploadUrlResponse, err error) {
	userID := ctxutil.GetUIDFromCtx(ctx)
	if userID == nil {
		return nil, errorx.New(errno.ErrUserInfoInvalidateCode, errorx.KV("msg", "User not authenticated"))
	}

	fileID, fileKey, url, err := i.ResumeDomainSVC.GetUploadURL(ctx, *userID, fileName, fileType)
	if err != nil {
		return nil, err
	}

	return &interviewAPI.ResumeUploadUrlResponse{
		UploadURL: url,
		FileKey:   fileKey,
		FileID:    fileID,
		Code:      0,
	}, nil
}

func (i *InterviewApplicationService) CreateResume(ctx context.Context, req *interviewAPI.ResumeMetaInfoRequest) (res *interviewAPI.ResumeMetaInfoResponse, err error) {
	userID := ctxutil.GetUIDFromCtx(ctx)
	if userID == nil {
		return nil, errorx.New(errno.ErrUserInfoInvalidateCode, errorx.KV("msg", "User not authenticated"))
	}

	resumeEntity, err := i.ResumeDomainSVC.Create(ctx, &service.ResumeCreateRequest{
		UserID:   *userID,
		ID:       req.FileID,
		FileKey:  req.FileKey,
		Filename: req.Filename,
		Filetype: req.Filetype,
		Filesize: req.Filesize,
	})
	if err != nil {
		return nil, err
	}

	// 异步发送Kafka消息，避免影响主流程性能
	go func() {
		// 构建消息结构
		msg := ResumeMsg{
			FileKey:  req.FileKey,
			FileID:   req.FileID,
			Filename: req.Filename,
			Filetype: req.Filetype,
			Filesize: req.Filesize,
			UserID:   *userID,
		}

		// 序列化消息
		msgJSON, err := json.Marshal(msg)
		if err != nil {
			logs.Errorf("Failed to marshal resume msg, userID: %d, fileID: %d, err: %v", *userID, req.FileID, err)
			return
		}

		// 发送消息，最多重试3次
		maxRetries := 3
		var sendErr error

		for retry := 0; retry < maxRetries; retry++ {
			sendErr = i.KafkaProducer.SendMessage(ctx, conf.Global.Kafka.ResumeTopic, []byte(req.FileKey), msgJSON)
			if sendErr == nil {
				logs.Infof("Successfully sent resume msg to Kafka, userID: %d, fileID: %d, fileKey: %s", *userID, req.FileID, req.FileKey)
				return
			}

			// 记录重试日志
			logs.Errorf("Failed to send resume msg to Kafka (attempt %d/%d), userID: %d, fileID: %d, err: %v",
				retry+1, maxRetries, *userID, req.FileID, sendErr)

			// 指数退避重试
			if retry < maxRetries-1 {
				time.Sleep(time.Duration(1<<uint(retry)) * 500 * time.Millisecond)
			}
		}

		// 所有重试都失败，记录最终错误
		logs.Errorf("All attempts to send resume msg to Kafka failed, userID: %d, fileID: %d, final err: %v",
			*userID, req.FileID, sendErr)
	}()

	return &interviewAPI.ResumeMetaInfoResponse{
		Data: resumeDo2UserTo(resumeEntity),
		Code: 0,
	}, nil
}

func resumeDo2UserTo(resumeDo *entity.Resume) *interviewAPI.ResumeInfo {
	return &interviewAPI.ResumeInfo{
		ID:       resumeDo.ID,
		UserID:   resumeDo.UserID,
		FileKey:  resumeDo.FileKey,
		Filename: resumeDo.Filename,
		Filetype: resumeDo.Filetype,
		Filesize: resumeDo.Filesize,
		Status:   resumeDo.Status,
		UploadAt: resumeDo.UploadAt,
	}
}
