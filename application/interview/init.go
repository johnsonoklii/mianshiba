package interview

import (
	"context"
	"mianshiba/domain/interview/repository"
	"mianshiba/domain/interview/service"
	"mianshiba/infra/contract/idgen"
	cmq "mianshiba/infra/contract/mq"
	"mianshiba/infra/contract/storage"

	"gorm.io/gorm"
)

func InitService(ctx context.Context, db *gorm.DB, idgen idgen.IDGenerator, minioClient storage.Storage, kafkaProducer cmq.KafkaProducer) *InterviewApplicationService {
	InterviewApplicationSVC.ResumeDomainSVC = service.NewResumeDomain(ctx, &service.ResumeComponents{
		OSSClient:  minioClient,
		IDGen:      idgen,
		ResumeRepo: repository.NewResumeRepo(db),
	})

	InterviewApplicationSVC.KafkaProducer = kafkaProducer

	return InterviewApplicationSVC
}
