package redis

import (
	"fmt"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"time"
)

type healthCheck struct {
	name string
}

func (c *healthCheck) Check() (string, error) {
	t, err := container.Resolve[IClient](c.name).Original().Time(fs.Context).Result()
	flog.ErrorIfExists(err)

	return fmt.Sprintf("Redis.%s => %s", c.name, t.Format(time.DateTime)), err
}
