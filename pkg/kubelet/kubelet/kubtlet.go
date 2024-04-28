package kubelet

import (
	"minik8s/pkg/kubelet/kconfig"
	"minik8s/pkg/kubelet/pod"

	// "minik8s/pkg/entity"
	// "minik8s/pkg/kubelet/event"
	// "minik8s/pkg/kubelet/lifecycle"
	"minik8s/pkg/kubelet/status"
	// "minik8s/pkg/kubelet/syncLoop"
)

type Kubelet struct {
	Config *kconfig.KubeletConfig
	// syncLoopHandler *syncLoop.SyncLoopHandler
	StatusManager status.StatusManager
	// plegManager     *lifecycle.PlegManager
	PodManager pod.PodManager

	// /* 用来接受syncLoop发送信号的通道，handler从该通道获取时间信号并响应处理 */
	// syncLoopChan chan *event.SyncLoopEventType

	// /* 用来更新pod信息的通道 */
	// podUpdateChan chan *entity.PodUpdateCmd
}

func (k *Kubelet) Run() {
	k.registerNode()
}

/* 需要到apiServer那边去注册相关信息，规定相关的数据结构和接口 */
func (k *Kubelet) registerNode() {
	_ = k.StatusManager.RegisterNode()
}
