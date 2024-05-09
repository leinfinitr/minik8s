// 描述：api server的配置

package config

import "strconv"

const (
	// APIServerLocalAddress api server的本地服务器地址
	APIServerLocalAddress = "127.0.0.1"
	// APIServerLocalPort api server的本地服务器端口
	APIServerLocalPort = 7000
)

func APIServerURL() string {
	return "http://" + APIServerLocalAddress + ":" + strconv.Itoa(APIServerLocalPort)
}

func NewAPIServerConfig() *APIServerConfig {
	return &APIServerConfig{
		APIServerIP:   APIServerLocalAddress,
		APIServerPort: APIServerLocalPort,
	}
}

type APIServerConfig struct {
	APIServerIP   string
	APIServerPort int
}

func (c *APIServerConfig) APIServerURL() string {
	return HttpSchema + c.APIServerIP + ":" + strconv.Itoa(c.APIServerPort)
}
