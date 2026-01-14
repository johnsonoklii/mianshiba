package handler

import (
	"context"
	"encoding/json"
	"mianshiba/application/interview"
	agentService "mianshiba/domain/agent/service"
	"mianshiba/domain/interview/event"
	cmq "mianshiba/infra/contract/mq"
	"mianshiba/pkg/logs"
	"time"
)

var ResumeHandlerSVC = &ResumeEventHandler{}

type ResumeEventHandler struct {
	ResumeAgentDomainSVC agentService.ResumeAgent
}

func (h *ResumeEventHandler) HandleResumeParseEvent(ctx context.Context, event *event.ResumeParseEvent) error {
	logs.Infof("Handling ResumeParseEvent: userID=%d, fileID=%d, fileKey=%s, filename=%s",
		event.UserID, event.FileID, event.FileKey, event.Filename)

	h.ResumeAgentDomainSVC.ParseResumeAndSave(ctx, &agentService.ParseResumeRequest{
		FileID:   event.FileID,
		UserID:   event.UserID,
		FileKey:  event.FileKey,
		Filename: event.Filename,
		Filetype: event.Filetype,
		Filesize: event.Filesize,
	})

	logs.Infof("Successfully handled ResumeCreatedEvent for file %s", event.FileKey)

	return nil
}

func ConvertToResumeParseDomainEvent(msg *cmq.KafkaMessage) (*event.ResumeParseEvent, error) {
	var resumeMsg interview.ResumeMsg
	if err := json.Unmarshal(msg.Value, &resumeMsg); err != nil {
		logs.Errorf("Failed to unmarshal resume msg: %v", err)
		return nil, err
	}

	// 创建领域事件
	domainEvent := &event.ResumeParseEvent{
		FileKey:   resumeMsg.FileKey,
		FileID:    resumeMsg.FileID,
		UserID:    resumeMsg.UserID,
		Filename:  resumeMsg.Filename,
		Filetype:  resumeMsg.Filetype,
		Filesize:  resumeMsg.Filesize,
		CreatedAt: time.Now(),
	}

	return domainEvent, nil
}
