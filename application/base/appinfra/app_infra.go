package appinfra

import (
	"context"
	"fmt"
	"mianshiba/infra/contract/cache"
	cmq "mianshiba/infra/contract/mq"
	"mianshiba/infra/contract/storage"
	"mianshiba/infra/impl/cache/redis"
	"mianshiba/infra/impl/idgen"
	mq "mianshiba/infra/impl/mq"
	"mianshiba/infra/impl/mysql"
	"mianshiba/infra/impl/storage/minio"

	"gorm.io/gorm"
)

type AppDependencies struct {
	DB            *gorm.DB
	CacheCli      cache.Cmdable
	IDGenSVC      idgen.IDGenerator
	MinIOClient   storage.Storage
	KafkaProducer cmq.KafkaProducer
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

	deps.MinIOClient, err = minio.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("init minio client failed, err=%w", err)
	}

	deps.KafkaProducer, err = mq.NewProducer(ctx)
	if err != nil {
		return nil, fmt.Errorf("init kafka producer failed, err=%w", err)
	}

	return deps, nil
}
