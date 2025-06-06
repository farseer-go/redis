package redis

import (
	"github.com/farseer-go/cache"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/modules"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{cache.Module{}}
}

func (module Module) PreInitialize() {
	container.Register(func() cache.ICache { return &cacheInRedis{} }, "redis")
}

func (module Module) Initialize() {
	redisConfigs := configure.GetSubNodes("Redis")
	for name, configString := range redisConfigs {
		Register(name, configString)
	}
}

// Register 注册Redis实例
func Register(name string, configString any) {
	config := configure.ParseString[redisConfig](configString.(string))
	if config.Server == "" {
		_ = flog.Error("Redis配置缺少Server节点")
		return
	}

	// 如果之前注册过，则先移除
	if container.IsRegister[IClient](name) {
		container.Remove[IClient](name)
	}
	container.Register(func() IClient { return newClient(config) }, name)

	// 如果之前注册过，则先移除
	if container.IsRegister[core.IHealthCheck]("redis_" + name) {
		container.Remove[core.IHealthCheck]("redis_" + name)
	}
	// 注册健康检查
	container.RegisterInstance[core.IHealthCheck](&healthCheck{name: name}, "redis_"+name)
}
