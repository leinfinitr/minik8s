package kubelet

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/config"
	"minik8s/pkg/kubelet/pod"
	"minik8s/pkg/kubelet/status"
	"minik8s/tools/host"
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
	k.registerKubeletAPI()
	KubeletIP, _ := host.GetHostIP()
	_ = k.KubeletAPIRouter.Run(KubeletIP + ":" + fmt.Sprint(config.KubeletAPIPort))
}

// registerNode 注册节点
func (k *Kubelet) registerNode() {
	_ = k.StatusManager.RegisterNode()
}

// registerKubeletAPI 注册kubelet的API
func (k *Kubelet) registerKubeletAPI() {
	k.KubeletAPIRouter.POST(config.PodsURI, pod.AddPod)
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
