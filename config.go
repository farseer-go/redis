package redis

// 定义redis配置结构
type redisConfig struct {
	Server          string
	DB              int
	Password        string
	ConnectTimeout  int
	SyncTimeout     int
	ResponseTimeout int
}
