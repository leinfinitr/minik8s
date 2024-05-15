package serverless

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/config"
	"minik8s/pkg/serverless/handler"
	"minik8s/tools/log"
)

type ServerlessServer struct {
	// 服务器地址
	Address string
	// 服务器端口
	Port int
	// 转发请求
	Router *gin.Engine
}

// 方法-------------------------------------------------------------

// Run 启动ApiServer
func (a *ServerlessServer) Run() {
	a.Register()
	err := a.Router.Run(a.Address + ":" + fmt.Sprint(a.Port))
	if err != nil {
		log.ErrorLog("ServerlessServer Run: " + err.Error())
	}
}

// Register 注册路由
func (a *ServerlessServer) Register() {
	// 创建Serverless环境
	a.Router.POST(config.ServerlessURI, handler.CreateServerless)
	// 获取所有的Serverless Function
	//a.Router.GET(config.ServerlessURI, handler.GetServerless)
	// 创建Serverless Function
	//a.Router.POST(config.ServerlessURI, handler.CreateServerless)
}

// 函数-------------------------------------------------------------

// NewServerlessServer 创建一个新的ServerlessServer
func NewServerlessServer() *ServerlessServer {
	return &ServerlessServer{
		Address: config.ServerlessAddress,
		Port:    config.ServerlessPort,
		Router:  gin.Default(),
	}
}
