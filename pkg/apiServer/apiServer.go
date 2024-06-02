// 描述：apiServer包实现了一个简单的HTTP服务器，用于处理API请求。
// 参考：gin框架 - https://www.topgoer.com/gin%E6%A1%86%E6%9E%B6/

package apiServer

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiObject"
	etcdclient "minik8s/pkg/apiServer/etcdClient"
	"minik8s/pkg/apiServer/handlers"
	"minik8s/pkg/config"
	httprequest "minik8s/tools/httpRequest"
	"minik8s/tools/log"
	"strings"
	"time"

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
	log.InfoLog("ApiServer is running on " + a.Address + ":" + fmt.Sprint(a.Port))

	go func() {
		a.Register()
		err := a.Router.Run(a.Address + ":" + fmt.Sprint(a.Port))
		if err != nil {
			panic(err)
		}
	}()

	// 开辟一个协程，用于定时扫描所有node的状态
	ScanNodeStatus()
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
	a.Router.PUT(config.NodeStatusURI, handlers.PingNodeStatus)

	// 获取指定Pod
	a.Router.GET(config.PodURI, handlers.GetPod)
	// 更新Pod
	a.Router.PUT(config.PodURI, handlers.UpdatePod)
	// 部分更新Pod
	a.Router.PATCH(config.PodURI, handlers.UpdatePod)
	// 删除指定Pod
	a.Router.DELETE(config.PodURI, handlers.DeletePod)

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
	// 获取全局所有Pod
	a.Router.GET(config.PodsGlobalURI, handlers.GetGlobalPods)

	// 接受kubeproxy的注册
	a.Router.POST(config.ProxyStatusURI, handlers.RegisterProxy)
	// 获取指定Service
	a.Router.GET(config.ServiceURI, handlers.GetService)
	// 更新指定Service
	a.Router.PUT(config.ServiceURI, handlers.PutService)
	// 删除制定Service
	a.Router.DELETE(config.ServiceURI, handlers.DeleteService)
	// 获取所有的services
	a.Router.GET(config.ServicesURI, handlers.GetServices)

	// 获取所有ReplicaSets
	a.Router.GET(config.ReplicaSetsURI, handlers.GetReplicaSets)
	// 获取全局所有ReplicaSets
	a.Router.GET(config.GlobalReplicaSetsURI, handlers.GetGlobalReplicaSets)
	// 获取指定ReplicaSet
	a.Router.GET(config.ReplicaSetURI, handlers.GetReplicaSet)
	//获取指定ReplicaSet的状态
	a.Router.GET(config.ReplicaSetStatusURI, handlers.GetReplicaSetStatus)
	//更新指定ReplicaSet的状态
	a.Router.POST(config.ReplicaSetStatusURI, handlers.UpdateReplicaSetStatus)
	//创建ReplicaSet
	a.Router.POST(config.ReplicaSetsURI, handlers.AddReplicaSet)
	//更新指定ReplicaSet
	a.Router.PUT(config.ReplicaSetURI, handlers.UpdateReplicaSet)
	//删除指定ReplicaSet
	a.Router.DELETE(config.ReplicaSetURI, handlers.DeleteReplicaSet)

	// 获取全局所有HPAs
	a.Router.GET(config.GlobalHpaURI, handlers.GetGlobalHPAs)
	// 获取所有HPAs
	a.Router.GET(config.HpasURI, handlers.GetHPAs)
	// 获取指定HPA
	a.Router.GET(config.HpaURI, handlers.GetHPA)
	// 创建指定HPA
	a.Router.POST(config.HpasURI, handlers.AddHPA)
	// 删除指定HPA
	a.Router.DELETE(config.HpaURI, handlers.DeleteHPA)
	// 更新指定HPA状态
	a.Router.PUT(config.HpaStatusURI, handlers.UpdateHPAStatus)

	// 获取所有DNSs
	a.Router.GET(config.DNSsURI, handlers.GetDNSs)
	// 获取指定DNS
	a.Router.GET(config.DNSURI, handlers.GetDNS)
	// 添加DNS
	a.Router.POST(config.DNSsURI, handlers.AddDNS)
	// 删除DNS
	a.Router.DELETE(config.DNSURI, handlers.DeleteDNS)
	// 删除所有DNS
	a.Router.DELETE(config.DnsRequestURI, handlers.DeleteDnsRequest)

	// 获取全局所有DNSs
	a.Router.GET(config.GlobalDnsRequestURI, handlers.GetGlobalDnsRequests)

	// 创建持久化卷
	a.Router.POST(config.PersistentVolumeURI, handlers.CreatePV)

	// 创建持久化卷声明
	a.Router.POST(config.PersistentVolumeClaimURI, handlers.CreatePVC)

	// 首次注册节点
	a.Router.PUT(config.MonitorNodeURL, handlers.RegisterNodeMonitor)
	// 节点失联后，删除相关配置
	a.Router.DELETE(config.MonitorNodeURL, handlers.DeleteNodeMonitor)
	// 注册暴露自定义指标的pod
	a.Router.PUT(config.MonitorPodURL, handlers.RegisterPodMonitor)
	// 删除暴露自定义指标的pod
	a.Router.DELETE(config.MonitorPodURL, handlers.DeletePodMonitor)

}

func ScanNodeStatus() {
	for {
		// 获取所有节点
		res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdNodePrefix)
		if err != nil {
			log.WarnLog("ScanNodeStatus: " + err.Error())
		}

		for _, v := range res {
			var node apiObject.Node
			err = json.Unmarshal([]byte(v), &node)
			if err != nil {
				log.WarnLog("ScanNodeStatus: " + err.Error())
			}
			url := config.APIServerURL() + config.NodeStatusURI
			url = strings.Replace(url, config.NameReplace, node.Metadata.Name, 1)
			httprequest.PutObjMsg(url, node)
		}

		time.Sleep(10 * time.Second)
	}

}

// NewApiServer 使用配置文件创建并返回一个新的ApiServer
func NewApiServer() *ApiServer {
	return &ApiServer{
		Address: config.APIServerLocalAddress,
		Port:    config.APIServerLocalPort,
		Router:  gin.New(),
	}
}
