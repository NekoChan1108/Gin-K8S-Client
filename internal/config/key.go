package config

// Key用于约束配置属性

type KeyName string

const (
	// ServerHost 服务IP
	ServerHost KeyName = "server_host"
	// ServerPort 监听端口
	ServerPort KeyName = "server_port"
	// ServerName 服务名
	ServerName KeyName = "server_name"
)
