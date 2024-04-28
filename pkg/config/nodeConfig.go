package config

const (
	// LocalNodeAddress 在本地运行的Node的地址
	LocalNodeAddress = "127.0.0.2"
	// KubeletAPIPort kubelet server 与 apiServer 通信的端口
	KubeletAPIPort = 10250
	// HealthPort 通过访问该端口可以判断 kubelet 是否正常工作
	HealthPort = 10248
	// cAdvisorPort cAdvisor的端口，通过该端口可以获取到该节点的环境信息以及 node 上运行的容器状态等内容
	cAdvisorPort = 4194
	// ReadOnlyPort 提供了 pod 和 node 的信息，接口以只读形式暴露出去，访问该端口不需要认证和鉴权
	ReadOnlyPort = 10255
)
