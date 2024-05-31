package config

import "strconv"

const (
	// NFSServer nfs服务器地址
	NFSServer = "192.168.1.12"

	// PVServerPath nfs服务器路径
	PVServerPath = "/pvserver"
	// PVClientPath nfs客户端路径
	PVClientPath = "/pvclient"

	// PVServerAddress 本地服务器地址
	PVServerAddress = "127.0.0.1"
	// PVServerPort 本地服务器端口
	PVServerPort = 7002
)

func PVServerURL() string {
	return HttpSchema + PVServerAddress + ":" + strconv.Itoa(PVServerPort)
}
