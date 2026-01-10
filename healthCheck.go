package redis

import (
	"context"
	"fmt"

	"github.com/farseer-go/fs/container"
)

type healthCheck struct {
	name string
}

func (c *healthCheck) Check() (string, error) {
	t, err := container.Resolve[IClient](c.name).Original().Time(context.Background()).Result()
	return fmt.Sprintf("Redis.%s => %s", c.name, t.Format("2006-01-02 15:04:05")), err
}
