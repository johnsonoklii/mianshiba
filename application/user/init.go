package user

import (
	"context"
	"mianshiba/domain/user/repository"
	"mianshiba/domain/user/service"
	"mianshiba/infra/contract/cache"
	"mianshiba/infra/contract/idgen"

	"gorm.io/gorm"
)

func InitService(ctx context.Context, db *gorm.DB, cacheCli cache.Cmdable, idgen idgen.IDGenerator) *UserApplicationService {
	UserApplicationSVC.UserDomainSVC = service.NewUserDomain(ctx, &service.UserComponents{
		CacheCli: cacheCli,
		IDGen:    idgen,
		UserRepo: repository.NewUserRepo(db),
	})
	UserApplicationSVC.UserModelDomainSVC = service.NewUserModelDomain(ctx, &service.UserModelComponents{
		CacheCli:      cacheCli,
		IDGen:         idgen,
		UserModelRepo: repository.NewUserModelRepo(db),
	})

	return UserApplicationSVC
}
