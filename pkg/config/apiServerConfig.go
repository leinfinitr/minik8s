// 描述：ip配置

package config

import "strconv"

const (
	// LocalServerAddress api server的本地服务器地址
	LocalServerAddress = "127.0.0.1"
	// LocalServerPort api server的本地服务器端口
	LocalServerPort = 7000

	HttpSchema = "http://"
)

func APIServerUrl() string {
	return HttpSchema + LocalServerAddress + ":" + strconv.Itoa(LocalServerPort)
}