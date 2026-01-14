package interview

import (
	"mianshiba/domain/interview/service"
	mq "mianshiba/infra/contract/mq"
)

var InterviewApplicationSVC = &InterviewApplicationService{}

type InterviewApplicationService struct {
	ResumeDomainSVC service.Resume
	KafkaProducer   mq.KafkaProducer
}
