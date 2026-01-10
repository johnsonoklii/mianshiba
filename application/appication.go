package application

import (
	"context"
	"fmt"
	"mianshiba/application/base/appinfra"
	"mianshiba/application/user"
)

type basicServices struct {
	infra   *appinfra.AppDependencies
	userSVC *user.UserApplicationService
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

	return &basicServices{
		infra:   infra,
		userSVC: userSVC,
	}, nil
}
