package main

import (
	"github.com/gin-gonic/gin"
	"minik8s/pkg/apiServer"
)

func main() {
	// 设置gin的运行模式
	gin.SetMode(gin.ReleaseMode)
	// 创建并运行一个新的ApiServer
	server := apiServer.NewApiServer()
	server.Run()
}
