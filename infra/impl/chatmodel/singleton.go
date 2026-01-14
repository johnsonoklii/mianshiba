package chatmodel

import (
	"mianshiba/infra/contract/chatmodel"
	"sync"
)

var (
	once             sync.Once
	singletonFactory chatmodel.Factory
)

func InitSingletonFactory(factory chatmodel.Factory) {
	once.Do(func() {
		singletonFactory = factory
	})
}

func GetSingletonFactory() chatmodel.Factory {
	return singletonFactory
}
