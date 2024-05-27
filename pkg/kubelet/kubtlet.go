package kubelet

import (
	"fmt"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/pkg/kubelet/pod"
	"minik8s/tools/host"
	"minik8s/tools/log"
	"minik8s/tools/netRequest"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Kubelet struct {
	// ApiServerConfig 存储apiServer的配置信息，用于和apiServer进行通信
	ApiServerConfig config.APIServerConfig

	// Iptables 存储每个pod的ip与UUID
	Iptables map[string]string

	// PodManager 用来管理Pod的信息
	PodManager pod.PodManager

	// KubeletAPIRouter 用来处理kubelet的请求
	KubeletAPIRouter *gin.Engine

	// 用来存储node的信息
	node *apiObject.Node

	//plegManager   *lifecycle.PlegManager

	// 用来接受syncLoop发送信号的通道，handler从该通道获取时间信号并响应处理
	//syncLoopChan chan *event.SyncLoopEventType

	// 用来更新pod信息的通道
	//podUpdateChan chan *entity.PodUpdateCmd
}

func (k *Kubelet) Run() {
	// 用于接受并转发来自与apiServer通信端口的请求
	go func() {
		k.registerKubeletAPI()
		KubeletIP, _ := host.GetHostIP()
		log.InfoLog("Listening and serving HTTP on " + KubeletIP + ":" + fmt.Sprint(config.KubeletAPIPort))
		_ = k.KubeletAPIRouter.Run(KubeletIP + ":" + fmt.Sprint(config.KubeletAPIPort))
	}()

	// 注册node
	k.registerNode()

	// 定时扫描pod的状态并进行相应的处理
	pod.ScanPodStatus()

}

// RegisterNode 在kubelet刚开始创建时，需要到apiServer的work node去注册
//
//	通过发送POST请求的方式去注册，默认API："/api/v1/nodes"
func (k *Kubelet) registerNode() bool {

	if k.node == nil {
		k.buildNode()
	}

	for {
		// 一致尝试注册直到成功为止
		url := k.ApiServerConfig.APIServerURL() + config.NodesURI

		statusCode, _, _ := netRequest.PostRequestByTarget(url, k.node)

		if statusCode != config.HttpSuccessCode {
			log.ErrorLog("register node failed")
		} else {
			log.InfoLog("register node success")
			return true
		}
		time.Sleep(15 * time.Second)
	}

}

func (k *Kubelet) buildNode() {
	// 注册所需的参数
	HostName, _ := host.GetHostname()
	HostIP, _ := host.GetHostIP()

	// 获取主机的内存大小
	capacity := make(map[string]string)
	totalMemory, _ := host.GetTotalMemory()
	capacity["memory"] = strconv.FormatUint(totalMemory, 10)

	// 获取主机的内存和CPU使用率
	allocatable := make(map[string]string)
	MemoryUsage, _ := host.GetMemoryUsageRate()
	CPUUsage, _ := host.GetCPULoad()
	allocatable["memory"] = strconv.FormatFloat(MemoryUsage, 'f', -1, 64)
	allocatable["cpu"] = strconv.FormatFloat(CPUUsage[0], 'f', -1, 64)

	k.node = &apiObject.Node{
		TypeMeta: apiObject.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		Metadata: apiObject.ObjectMeta{
			Name:        HostName,
			Namespace:   "", // 该字段为空
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			UUID:        "", //	由API Server生成
		},
		Spec: apiObject.NodeSpec{
			PodCIDR:       "", // 未使用
			ProviderID:    "", // 未使用
			Unschedulable: true,
		},
		Status: apiObject.NodeStatus{
			Capacity:    capacity,
			Allocatable: allocatable,
			Phase:       "running",
			Conditions: []apiObject.NodeCondition{
				{
					Type:   "Ready", // Ready: kubelet准备好接受Pod
					Status: "True",
				},
			},
			Addresses: []apiObject.NodeAddress{
				{
					Type:    "InternalIP",
					Address: HostIP,
				},
			},
		},
	}

}

func (k *Kubelet) UpdateNodeStatusInternal() {
	log.InfoLog("UpdateNodeStatus")
	// 获取主机的内存和CPU使用率
	allocatable := make(map[string]string)
	MemoryUsage, _ := host.GetMemoryUsageRate()
	CPUUsage, _ := host.GetCPULoad()
	allocatable["memory"] = strconv.FormatFloat(MemoryUsage, 'f', -1, 64)
	allocatable["cpu"] = strconv.FormatFloat(CPUUsage[0], 'f', -1, 64)

	nodeStatus := &apiObject.NodeStatus{
		Allocatable: allocatable,
		Conditions: []apiObject.NodeCondition{
			{
				Type:   "Ready",
				Status: "True",
			},
		},
	}

	k.node.Status = *nodeStatus

}

// registerKubeletAPI 注册kubelet的API
func (k *Kubelet) registerKubeletAPI() {
	log.DebugLog("register kubelet API")
	// 该部分实现与 apiServer 中保持一致，每个方法的作用也参考 pkg/apiServer/apiServer.go 中的注释

	// 获取所有节点
	// k.KubeletAPIRouter.GET(config.NodesURI, handlers.GetNodes)
	// 创建节点
	// k.KubeletAPIRouter.POST(config.NodesURI, handlers.CreateNode)
	// 删除所有节点
	// k.KubeletAPIRouter.DELETE(config.NodesURI, handlers.DeleteNodes)

	// 获取指定节点
	// k.KubeletAPIRouter.GET(config.NodeURI, handlers.GetNode)
	// 更新指定节点
	// k.KubeletAPIRouter.PUT(config.NodeURI, handlers.UpdateNode)
	// 部分更新指定节点
	// k.KubeletAPIRouter.PATCH(config.NodeURI, handlers.UpdateNode)
	// 删除指定节点
	// k.KubeletAPIRouter.DELETE(config.NodeURI, handlers.DeleteNode)

	// 获取指定节点的状态
	// k.KubeletAPIRouter.GET(config.NodeStatusURI, handlers.GetNodeStatus)
	// 更新指定节点的状态
	// k.KubeletAPIRouter.PUT(config.NodeStatusURI, handlers.UpdateNodeStatus)
	// 部分更新指定节点的状态
	k.KubeletAPIRouter.GET(config.NodeStatusURI, k.UpdateNodeStatus)

	// 获取指定Pod
	// k.KubeletAPIRouter.GET(config.PodURI, handlers.GetPod)
	// 更新Pod
	k.KubeletAPIRouter.PUT(config.PodURI, pod.UpdatePod)
	// 部分更新Pod
	// k.KubeletAPIRouter.PATCH(config.PodURI, handlers.UpdatePod)
	// 删除指定Pod
	k.KubeletAPIRouter.DELETE(config.PodURI, pod.DeletePod)
	// kubelet挂掉了，apiServer用来同步Pod的信息
	k.KubeletAPIRouter.POST(config.PodsSyncURI, pod.SyncPods)

	// 获取指定Pod的状态
	k.KubeletAPIRouter.GET(config.PodStatusURI, pod.GetPodStatus)

	// 执行指定Pod和container的命令
	k.KubeletAPIRouter.GET(config.PodExecURI, pod.ExecPodContainer)

	// 获取所有Pod
	k.KubeletAPIRouter.GET(config.PodsURI, pod.GetPods)
	// 创建Pod
	k.KubeletAPIRouter.POST(config.PodsURI, pod.CreatePod)
	// 删除所有Pod
	// k.KubeletAPIRouter.DELETE(config.PodsURI, handlers.DeletePods)
}

// NewKubelet 创建一个新的Kubelet
func NewKubelet() *Kubelet {
	return &Kubelet{
		ApiServerConfig:  *config.NewAPIServerConfig(),
		PodManager:       pod.GetPodManager(),
		KubeletAPIRouter: gin.Default(),
		node:             nil,
	}
}

func (k *Kubelet) UpdateNodeStatus(c *gin.Context) {
	log.InfoLog("UpdateNodeStatus")
	if k.node == nil {
		k.buildNode()
	} else {
		k.UpdateNodeStatusInternal()
	}
	c.JSON(config.HttpSuccessCode, *k.node)
}
