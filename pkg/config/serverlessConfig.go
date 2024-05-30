package config

import "strconv"

const (
	// ServerlessAddress serverless的地址
	ServerlessAddress = "127.0.0.1"
	// ServerlessPort serverless的端口
	ServerlessPort = 7001
)

func ServerlessURL() string {
	return HttpSchema + ServerlessAddress + ":" + strconv.Itoa(ServerlessPort)
}
