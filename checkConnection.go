package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/trace"
)

// 确保实现了IConnectionChecker接口
var _ core.IConnectionChecker = (*connectionChecker)(nil)

type connectionChecker struct{}

// Check 检查连接字符串是否能成功连接到Redis
// 实现IConnectionChecker接口
func (c *connectionChecker) Check(configString string) (bool, error) {
	// 取消链路
	trace.Manager().Ignore()

	if configString == "" {
		return false, fmt.Errorf("连接字符串不能为空")
	}

	// 解析配置字符串
	config := configure.ParseString[redisConfig](configString)

	if config.Server == "" {
		return false, fmt.Errorf("Server配置不正确：%s", configString)
	}

	// 创建Redis客户端（使用newClient）
	connectTimeout := time.Duration(config.ConnectTimeout) * time.Millisecond
	if connectTimeout == 0 {
		connectTimeout = 5 * time.Second
	}

	cli := newClient(config)
	if cli == nil {
		return false, fmt.Errorf("创建Redis客户端失败")
	}
	rdb := cli.Original()
	if rdb == nil {
		return false, fmt.Errorf("获取原生Redis客户端失败")
	}
	defer rdb.Close()

	// Ping测试连接
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return false, fmt.Errorf("连接Redis失败：%s", err.Error())
	}

	return true, nil
}

// CheckConnection 检查连接字符串是否能成功连接到Redis
// configString 格式参考：
// Server=127.0.0.1:6379,DB=0,Password=,ConnectTimeout=5000,SyncTimeout=5000,ResponseTimeout=5000
// 返回值：(成功, 错误信息)
// CheckWithTimeout 带超时时间的连接检查
// timeout 为0时使用默认的10秒超时
func (c *connectionChecker) CheckWithTimeout(configString string, timeout time.Duration) (bool, error) {
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	type result struct {
		success bool
		err     error
	}
	resultChan := make(chan result, 1)

	go func() {
		success, err := c.Check(configString)
		resultChan <- result{success: success, err: err}
	}()

	select {
	case <-ctx.Done():
		return false, fmt.Errorf("连接检查超时，超时时间：%v", timeout)
	case res := <-resultChan:
		return res.success, res.err
	}
}
