package kubelet

import (
	"minik8s/pkg/kubelet/kconfig"
	// "minik8s/pkg/entity"
	// "minik8s/pkg/kubelet/event"
	// "minik8s/pkg/kubelet/lifecycle"
	// "minik8s/pkg/kubelet/status"
	// "minik8s/pkg/kubelet/syncLoop"
)

/* Kubelet数据结构主要集成了该组件的其他模块信息 */
type Kubelet struct {
	config *kconfig.KubeletConfig
	// syncLoopHandler *syncLoop.SyncLoopHandler
	// statusManager   *status.StatusManeger
	// plegManager     *lifecycle.PlegManager

	// /* 用来接受syncLoop发送信号的通道，handler从该通道获取时间信号并响应处理 */
	// syncLoopChan chan *event.SyncLoopEventType

	// /* 用来更新pod信息的通道 */
	// podUpdateChan chan *entity.PodUpdateCmd
}

func Run() {

}
