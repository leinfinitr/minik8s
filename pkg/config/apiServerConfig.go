package config

import "strconv"

const (
	// APIServerLocalAddress api server的本地服务器地址
	APIServerLocalAddress = "192.168.1.12"
	// APIServerLocalPort api server的本地服务器端口
	APIServerLocalPort = 7000
)

type APIServerConfig struct {
	APIServerIP   string
	APIServerPort int
}

func (c *APIServerConfig) APIServerURL() string {
	return HttpSchema + c.APIServerIP + ":" + strconv.Itoa(c.APIServerPort)
}

func APIServerURL() string {
	return HttpSchema + APIServerLocalAddress + ":" + strconv.Itoa(APIServerLocalPort)
}

func NewAPIServerConfig() *APIServerConfig {
	return &APIServerConfig{
		APIServerIP:   APIServerLocalAddress,
		APIServerPort: APIServerLocalPort,
	}
}
