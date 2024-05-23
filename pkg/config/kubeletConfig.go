package config

const (
	// KubeletLocalIP kubelet server 的地址
	KubeletLocalURLPrefix = "http://"

	// KubeletAPIPort kubelet server 与 apiServer 通信的端口
	KubeletAPIPort = 10250
	// KubeletHealthPort 通过访问该端口可以判断 kubelet 是否正常工作
	KubeletHealthPort = 10248
	// cAdvisorPort cAdvisor的端口，通过该端口可以获取到该节点的环境信息以及 node 上运行的容器状态等内容
	cAdvisorPort = 4194
	// ReadOnlyPort 提供了 pod 和 node 的信息，接口以只读形式暴露出去，访问该端口不需要认证和鉴权
	ReadOnlyPort = 10255
)

// grpc request
const (
	ContainerRuntimeEndpoint = "unix:///run/containerd/containerd.sock"
	ImageRuntimeEndpoint     = "unix:///run/containerd/containerd.sock"

	MaxMsgSize = 1024 * 1024 * 16
)

// defaultUnixEndpoint = "unix:///tmp/kubelet_remote_%v.sock"
