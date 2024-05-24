// 描述：apiServer包实现了一个简单的HTTP服务器，用于处理API请求。
// 参考：gin框架 - https://www.topgoer.com/gin%E6%A1%86%E6%9E%B6/

package apiServer

import (
	"fmt"
	"minik8s/pkg/apiServer/handlers"
	"minik8s/pkg/config"

	"github.com/gin-gonic/gin"
)

type ApiServer struct {
	// 服务器地址
	Address string
	// 服务器端口
	Port int
	// 转发请求
	Router *gin.Engine
}

// 方法-------------------------------------------------------------

// Run 启动ApiServer
func (a *ApiServer) Run() {
	a.Register()
	err := a.Router.Run(a.Address + ":" + fmt.Sprint(a.Port))
	if err != nil {
		panic(err)
	}
}

// Register 注册路由
//
//	Node - https://kubernetes.io/zh-cn/docs/reference/kubernetes-api/cluster-resources/node-v1/#Operations
//	Pod - https://kubernetes.io/zh-cn/docs/reference/kubernetes-api/cluster-resources/pod-v1/#Operations
func (a *ApiServer) Register() {
	// GET: 查询
	// POST: 创建
	// PUT: 更新
	// DELETE: 删除
	// PATCH: 更新部分资源

	// 获取所有节点
	a.Router.GET(config.NodesURI, handlers.GetNodes)
	// 创建节点
	a.Router.POST(config.NodesURI, handlers.CreateNode)
	// 删除所有节点
	a.Router.DELETE(config.NodesURI, handlers.DeleteNodes)

	// 获取指定节点
	a.Router.GET(config.NodeURI, handlers.GetNode)
	// 更新指定节点
	a.Router.PUT(config.NodeURI, handlers.UpdateNode)
	// 部分更新指定节点
	a.Router.PATCH(config.NodeURI, handlers.UpdateNode)
	// 删除指定节点
	a.Router.DELETE(config.NodeURI, handlers.DeleteNode)

	// 获取指定节点的状态
	a.Router.GET(config.NodeStatusURI, handlers.GetNodeStatus)
	// 更新指定节点的状态
	a.Router.PUT(config.NodeStatusURI, handlers.UpdateNodeStatus)
	// 部分更新指定节点的状态
	a.Router.PATCH(config.NodeStatusURI, handlers.UpdateNodeStatus)

	// 获取指定Pod
	a.Router.GET(config.PodURI, handlers.GetPod)
	// 更新Pod
	a.Router.PUT(config.PodURI, handlers.UpdatePod)
	// 部分更新Pod
	a.Router.PATCH(config.PodURI, handlers.UpdatePod)
	// 删除指定Pod
	a.Router.DELETE(config.PodURI, handlers.DeletePod)

	// 获取指定Pod的EphemeralContainers
	a.Router.GET(config.PodEphemeralContainersURI, handlers.GetPodEphemeralContainers)
	// 更新Pod的EphemeralContainers
	a.Router.PUT(config.PodEphemeralContainersURI, handlers.UpdatePodEphemeralContainers)
	// 部分更新Pod的EphemeralContainers
	a.Router.PATCH(config.PodEphemeralContainersURI, handlers.UpdatePodEphemeralContainers)

	// 获取指定Pod的日志
	a.Router.GET(config.PodLogURI, handlers.GetPodLog)

	// 获取指定Pod的状态
	a.Router.GET(config.PodStatusURI, handlers.GetPodStatus)
	// 更新Pod的状态
	a.Router.PUT(config.PodStatusURI, handlers.UpdatePodStatus)

	// 执行指定Pod和container的命令
	a.Router.GET(config.PodExecURI, handlers.ExecPod)

	// 获取所有Pod
	a.Router.GET(config.PodsURI, handlers.GetPods)
	// 创建Pod
	a.Router.POST(config.PodsURI, handlers.CreatePod)
	// 删除所有Pod
	a.Router.DELETE(config.PodsURI, handlers.DeletePods)

	// 接受kubeproxy的心跳
	a.Router.PUT(config.ProxiesStatusURI, handlers.UpdateProxyStatus)
	// 获取指定Service
	a.Router.GET(config.ServiceURI, handlers.GetService)
	// 更新指定Service
	a.Router.PUT(config.ServiceURI, handlers.PutService)
	// 删除制定Service
	a.Router.DELETE(config.ServiceURI, handlers.GetService)

}

// 函数-------------------------------------------------------------

// NewApiServer 使用配置文件创建并返回一个新的ApiServer
func NewApiServer() *ApiServer {
	return &ApiServer{
		Address: config.APIServerLocalAddress,
		Port:    config.APIServerLocalPort,
		Router:  gin.Default(),
	}
}
