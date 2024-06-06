package kubelet

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/pkg/kubelet/pod"
	"minik8s/tools/host"
	"minik8s/tools/log"
	"minik8s/tools/netRequest"
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
func (k *Kubelet) registerNode() bool {
	if k.node == nil {
		k.buildNode()
	}

	// 一直尝试注册直到成功为止
	for {
		url := k.ApiServerConfig.APIServerURL() + config.NodesURI
		statusCode, _, _ := netRequest.PostRequestByTarget(url, *k.node)
		if statusCode != config.HttpSuccessCode {
			log.ErrorLog("register node failed")
		} else {
			log.InfoLog("register node success")
			return true
		}
		time.Sleep(15 * time.Second)
	}

}

// buildNode 构建node的信息
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
			Namespace:   "",
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			UUID:        "",
		},
		Spec: apiObject.NodeSpec{
			PodCIDR:       "",
			ProviderID:    "",
			Unschedulable: true,
		},
		Status: apiObject.NodeStatus{
			Capacity:    capacity,
			Allocatable: allocatable,
			Phase:       "running",
			Conditions: []apiObject.NodeCondition{
				{
					Type:   "Ready",
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

// UpdateNodeStatusInternal 更新node的状态
func (k *Kubelet) UpdateNodeStatusInternal() {
	log.DebugLog("UpdateNodeStatus")

	// 注册所需的参数
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

	nodeStatus := &apiObject.NodeStatus{
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
	}

	k.node.Status = *nodeStatus

}

// registerKubeletAPI 注册kubelet的API
func (k *Kubelet) registerKubeletAPI() {
	// 该部分接口实现与 apiServer 中保持一致
	log.DebugLog("register kubelet API")

	// 部分更新指定节点的状态
	k.KubeletAPIRouter.GET(config.NodeStatusURI, k.UpdateNodeStatus)

	// 更新Pod
	k.KubeletAPIRouter.PUT(config.PodURI, pod.UpdatePod)
	// 删除指定Pod
	k.KubeletAPIRouter.DELETE(config.PodURI, pod.DeletePod)
	// kubelet挂掉之后，apiServer用来同步Pod的信息
	k.KubeletAPIRouter.POST(config.PodsSyncURI, pod.SyncPods)

	// 获取指定Pod的状态
	k.KubeletAPIRouter.GET(config.PodStatusURI, pod.GetPodStatus)

	// 执行指定Pod和container的命令
	k.KubeletAPIRouter.POST(config.PodExecURI, pod.ExecPodContainer)

	// 获取所有Pod
	k.KubeletAPIRouter.GET(config.PodsURI, pod.GetPods)
	// 创建Pod
	k.KubeletAPIRouter.POST(config.PodsURI, pod.CreatePod)
}

// NewKubelet 创建一个新的Kubelet
func NewKubelet() *Kubelet {
	return &Kubelet{
		ApiServerConfig:  *config.NewAPIServerConfig(),
		PodManager:       pod.GetPodManager(),
		KubeletAPIRouter: gin.New(),
		node:             nil,
	}
}

// UpdateNodeStatus 更新node的状态
func (k *Kubelet) UpdateNodeStatus(c *gin.Context) {
	log.DebugLog("UpdateNodeStatus")
	if k.node == nil {
		k.buildNode()
	} else {
		k.UpdateNodeStatusInternal()
	}
	c.JSON(config.HttpSuccessCode, k.node.Status)
}
