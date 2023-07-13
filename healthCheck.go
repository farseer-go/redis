package redis

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
)

type healthCheck struct {
	name string
}

func (c *healthCheck) Check() (string, error) {
	_, err := container.Resolve[IClient](c.name).Original().Time(fs.Context).Result()
	flog.ErrorIfExists(err)
	return "Redis." + c.name, err
}
