package kconfig

import (
	"minik8s/pkg/ipconfig"
)

// TODO: 需要的API信息还没有完全确定，暂时先这样

/* Kubelet需要与API Server进行交互，所以记录API Server的相关信息 */
type KubeletConfig struct {
	apiServerIP             string
	apiServerPort           int
	apiServerCAFile         string
	apiServerClientCertFile string
}

func kubeletConfigDefault() *KubeletConfig {
	apiLocalIP := ipconfig.APILocalIP
	apiServerPort := ipconfig.APIServerPort

	return &KubeletConfig{
		apiServerIP:             apiLocalIP,
		apiServerPort:           apiServerPort,
		apiServerCAFile:         "",
		apiServerClientCertFile: "",
	}
}
