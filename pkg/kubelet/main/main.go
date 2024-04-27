package main

import (
	"minik8s/pkg/kubelet/kconfig"
	"minik8s/pkg/kubelet/kubelet"
	"minik8s/pkg/kubelet/pod"
	"minik8s/pkg/kubelet/status"
	// "minik8s/pkg/kubelet/syncLoop"
)

func main() {
	// apiServer那边的配置暂时采用默认值
	newKubeletConfig := kconfig.KubeletConfigDefault()

	kubeletInit, err := newKubeletInit(newKubeletConfig)
	if err != nil {
		panic(err)
	}

	kubeletInit.Run()
}

// 包括调用各种模块的初始化函数，并在Kubelet中组装注册
func newKubeletInit(config *kconfig.KubeletConfig) (*kubelet.Kubelet, error) {
	// 处理传进来的相关参数，配置config，调用run
	statusTmp, err := status.GetStatusManager(config.APIServerURL, config.APIServerIP)
	if err != nil {
		panic(err)
	}

	K := &kubelet.Kubelet{
		Config:        config,
		PodManager:    pod.GetPodManager(),
		StatusManager: statusTmp,
		// syncLoopHandler:	syncLoop.SyncLoopHandler,

	}
	return K, nil
}
