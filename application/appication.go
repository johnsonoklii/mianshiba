package application

import (
	"context"
	"fmt"
	"mianshiba/application/agent"
	"mianshiba/application/agent/handler"
	"mianshiba/application/base/appinfra"
	"mianshiba/application/interview"
	"mianshiba/application/user"
)

type basicServices struct {
	infra        *appinfra.AppDependencies
	userSVC      *user.UserApplicationService
	interviewSVC *interview.InterviewApplicationService
	agentHandler *handler.ResumeEventHandler
}

func Init(ctx context.Context) (err error) {
	infra, err := appinfra.Init(ctx)
	if err != nil {
		return err
	}

	_, err = initBasicServices(ctx, infra)
	if err != nil {
		return fmt.Errorf("Init - initBasicServices failed, err: %v", err)
	}

	return nil
}

// initBasicServices init basic services that only depends on infra.
func initBasicServices(ctx context.Context, infra *appinfra.AppDependencies) (*basicServices, error) {
	userSVC := user.InitService(ctx, infra.DB, infra.CacheCli, infra.IDGenSVC)
	interviewSVC := interview.InitService(ctx, infra.DB, infra.IDGenSVC, infra.MinIOClient, infra.KafkaProducer)
	agentHandler := agent.InitHandler(ctx, infra.DB, infra.MinIOClient)

	return &basicServices{
		infra:        infra,
		userSVC:      userSVC,
		interviewSVC: interviewSVC,
		agentHandler: agentHandler,
	}, nil
}
