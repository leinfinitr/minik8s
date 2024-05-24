package config

const (
	// KubeproxyLocalAddress kubeproxy的本地服务器地址
	KubeproxyLocalAddress   = "127.0.0.1"
	KubeproxyLocalUTLPrefix = "http://" + KubeproxyLocalAddress

	// KubeproxyLocalPort kubeproxy的本地服务器端口
	KubeproxyAPIPort = 10256
	// HealthPort 通过访问该端口可以判断 kubeproxy 是否正常工作
	KubeproyHealthPort = 10249
)
