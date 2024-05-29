package serverless

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"minik8s/pkg/config"
	"minik8s/pkg/serverless/handler"
	"minik8s/pkg/serverless/scale"
	"minik8s/tools/log"
)

type ServerlessServer struct {
	// 服务器地址
	Address string
	// 服务器端口
	Port int
	// 转发请求
	Router *gin.Engine
	// 自动扩容控制
	Scale scale.ScaleManagerImpl
}

// Run 启动ServerlessServer
func (s *ServerlessServer) Run() {
	s.Register()

	// 开启一个线程运行自动扩容控制
	go s.Scale.Run()

	// 主线程用于处理请求
	err := s.Router.Run(s.Address + ":" + fmt.Sprint(s.Port))
	if err != nil {
		log.ErrorLog("ServerlessServer Run: " + err.Error())
	}
}

// Register 注册路由
func (s *ServerlessServer) Register() {
	// 创建Serverless Function环境
	s.Router.POST(config.ServerlessURI, handler.CreateServerless)
	// 获取所有的Serverless Function
	s.Router.GET(config.ServerlessURI, handler.GetServerless)

	// 删除Serverless Function
	s.Router.DELETE(config.ServerlessFunctionURI, handler.DeleteServerless)
	// 更新Serverless Function
	s.Router.PUT(config.ServerlessFunctionURI, handler.UpdateServerlessFunction)

	// 运行Serverless Function
	s.Router.GET(config.ServerlessRunURI, handler.RunServerlessFunction)

	// 运行Serverless Workflow
	s.Router.POST(config.ServerlessWorkflowURI, handler.RunServerlessWorkflow)
}

// NewServerlessServer 创建一个新的ServerlessServer
func NewServerlessServer() *ServerlessServer {
	return &ServerlessServer{
		Address: config.ServerlessAddress,
		Port:    config.ServerlessPort,
		Router:  gin.Default(),
		Scale:   *scale.NewScaleManager(),
	}
}
