package main

// import (
// 	"minik8s/pkg/kubelet/kconfig"
// 	"minik8s/pkg/kubelet/kubelet"
// 	// "minik8s/pkg/kubelet/syncLoop"
// )

// func main() {
// 	// apiserver那边的配置暂时采用默认值
// 	newKubletConfig := kconfig.kubeletConfigDefault()

// 	kubelet, err := newKubeletInit(newKubletConfig)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// kubelet.Run()
// }

// // cobra 从kubectl那边进行解析，在这里不进行解析
// func newKubeletInit(config *kconfig.KubeletConfig) (*kubelet.Kubelet, error) {
// 	// 处理传进来的相关参数，配置config，调用run
// 	K := &kubelet.Kubelet{
// 		// config: 			config,
// 		// syncLoopHandler:	syncLoop.SyncLoopHandler,

// 	}
// 	return K, nil
// }
