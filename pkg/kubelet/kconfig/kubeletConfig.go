package kconfig

import (
	"minik8s/pkg/ipconfig"
)

// TODO: 需要的API信息还没有完全确定，暂时先这样

/* Kubelet需要与API Server进行交互，所以记录API Server的相关信息 */
type KubeletConfig struct {
	APIServerIP             string
	APIServerPort           int
	APIServerURL            string
	APIServerClientCertFile string
}

func KubeletConfigDefault() *KubeletConfig {
	apiLocalIP := ipconfig.APILocalIP
	apiServerPort := ipconfig.APIServerPort

	return &KubeletConfig{
		APIServerIP:             apiLocalIP,
		APIServerPort:           apiServerPort,
		APIServerURL:            "127.0.0.1:7000",
		APIServerClientCertFile: "",
	}
}
