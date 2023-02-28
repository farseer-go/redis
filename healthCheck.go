package redis

import (
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
)

type healthCheck struct {
	name string
}

func (c *healthCheck) Check() (string, error) {
	_, err := container.Resolve[IClient](c.name).Original().Time(fs.Context).Result()
	return "Redis." + c.name, err
}
