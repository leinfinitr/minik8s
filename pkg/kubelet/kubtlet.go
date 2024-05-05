package kubelet

import (
	"fmt"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/pkg/kubelet/pod"
	"minik8s/pkg/kubelet/status"
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
	// StatusManager 用来管理Node的状态信息
	StatusManager status.StatusManager
	// PodManager 用来管理Pod的信息
	PodManager pod.PodManager

	// KubeletAPIRouter 用来处理kubelet的请求
	KubeletAPIRouter *gin.Engine

	//plegManager   *lifecycle.PlegManager

	// 用来接受syncLoop发送信号的通道，handler从该通道获取时间信号并响应处理
	//syncLoopChan chan *event.SyncLoopEventType

	// 用来更新pod信息的通道
	//podUpdateChan chan *entity.PodUpdateCmd
}

func (k *Kubelet) Run() {
	k.registerNode()

	go func() {
		k.registerKubeletAPI()
		KubeletIP, _ := host.GetHostIP()
		_ = k.KubeletAPIRouter.Run(KubeletIP + ":" + fmt.Sprint(config.KubeletAPIPort))
	}()

	// kubelet的主线程用于发送心跳
	k.heartbeat()
}

// RegisterNode 在kubelet刚开始创建时，需要到apiServer的work node去注册
//
//	通过发送POST请求的方式去注册，默认API："/api/v1/nodes"
func (k *Kubelet) registerNode() bool {
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

	node := &apiObject.Node{
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

	url := k.ApiServerConfig.APIServerURL() + config.NodesURI

	statusCode, _, _ := netRequest.PostRequestByTarget(url, node)

	if statusCode != config.HttpSuccessCode {
		log.ErrorLog("register node failed")
		return false
	} else {
		log.InfoLog("register node success")
		return true
	}
}

// registerKubeletAPI 注册kubelet的API
func (k *Kubelet) registerKubeletAPI() {
	log.DebugLog("register kubelet API")
	// 该部分实现与 apiServer 中保持一致，每个方法的作用也参考 pkg/apiServer/apiServer.go 中的注释
	k.KubeletAPIRouter.DELETE(config.PodURI, pod.DeletePod)
	k.KubeletAPIRouter.POST(config.PodsURI, pod.AddPod)
}

// heartbeat 向apiServer发送心跳
func (k *Kubelet) heartbeat() {
	log.DebugLog("start heartbeat")
	// 每间隔 60s 发送一次心跳
	for {
		HostName, _ := host.GetHostname()
		url := k.ApiServerConfig.APIServerURL() + config.NodesURI + "/" + HostName + "/status"
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

		statusCode, _, _ := netRequest.PutRequestByTarget(url, nodeStatus)

		if statusCode != config.HttpSuccessCode {
			log.ErrorLog("heartbeat failed")
		} else {
			log.DebugLog("heartbeat success")
		}

		// 间隔60s
		time.Sleep(60 * time.Second)
	}
}

// NewKubelet 创建一个新的Kubelet
func NewKubelet() *Kubelet {
	return &Kubelet{
		ApiServerConfig:  *config.NewAPIServerConfig(),
		StatusManager:    status.GetStatusManager(config.APIServerURL()),
		PodManager:       pod.GetPodManager(),
		KubeletAPIRouter: gin.Default(),
	}
}
