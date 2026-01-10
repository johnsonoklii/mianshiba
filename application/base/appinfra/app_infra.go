package appinfra

import (
	"context"
	"fmt"
	"mianshiba/infra/contract/cache"
	"mianshiba/infra/impl/cache/redis"
	"mianshiba/infra/impl/idgen"
	"mianshiba/infra/impl/mysql"

	"gorm.io/gorm"
)

type AppDependencies struct {
	DB       *gorm.DB
	CacheCli cache.Cmdable
	IDGenSVC idgen.IDGenerator
}

func Init(ctx context.Context) (*AppDependencies, error) {
	deps := &AppDependencies{}
	var err error

	deps.DB, err = mysql.New()
	if err != nil {
		return nil, fmt.Errorf("init db failed, err=%w", err)
	}

	deps.CacheCli = redis.New()

	deps.IDGenSVC, err = idgen.New(deps.CacheCli)
	if err != nil {
		return nil, fmt.Errorf("init id gen svc failed, err=%w", err)
	}

	return deps, nil
}
