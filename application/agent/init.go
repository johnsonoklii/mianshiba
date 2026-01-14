package agent

import (
	"context"
	"mianshiba/application/agent/handler"
	agentService "mianshiba/domain/agent/service"
	"mianshiba/domain/interview/repository"
	"mianshiba/infra/contract/storage"

	"gorm.io/gorm"
)

func InitHandler(ctx context.Context, db *gorm.DB, minioClient storage.Storage) *handler.ResumeEventHandler {
	handler.ResumeHandlerSVC.ResumeAgentDomainSVC = agentService.NewResumeAgent(&agentService.ResumeAgentComponents{
		OSSClient:  minioClient,
		ResumeRepo: repository.NewResumeRepo(db),
	})

	return handler.ResumeHandlerSVC
}
