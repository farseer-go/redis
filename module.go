package redis

import (
	"fmt"
	"github.com/farseer-go/cache"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/modules"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{cache.Module{}}
}

func (module Module) PreInitialize() {
	container.Use[cache.ICache](func() cache.ICache { return cacheInRedis{} }).Name("redis").Register()
	cache := container.ResolveName[cache.ICache]("redis")
	fmt.Print(cache)
}

func (module Module) Initialize() {
}

func (module Module) PostInitialize() {
}

func (module Module) Shutdown() {
}
