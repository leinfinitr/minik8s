// 描述: 容器相关的 API 对象
// 参考：https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/api/core/v1/types.go

package apiObject

type ContainerStatus string

const (
	Created  ContainerStatus = "created"
	Running  ContainerStatus = "running"
	Paused   ContainerStatus = "paused"
	Restart  ContainerStatus = "restarting"
	Removing ContainerStatus = "removing"
	Exited   ContainerStatus = "exited"
	Dead     ContainerStatus = "dead"
)

type Container struct {
	// Name: 容器的名称
	Name string `json:"name" yaml:"name"`
	// Image: 容器的镜像
	Image string `json:"image" yaml:"image"`
	// Command: 容器的命令
	Command []string `json:"command" yaml:"command"`
	// Args: 容器的命令行参数
	Args []string `json:"args" yaml:"args"`

	// WorkingDir: 容器的工作目录
	WorkingDir string `json:"workingDir" yaml:"workingDir"`
	// Ports: 容器的端口
	Ports []ContainerPort `json:"ports" yaml:"ports"`
	// Env: 容器的环境变量
	Env []EnvVar `json:"env" yaml:"env"`
	// Resources: 容器的资源限制
	Resources ResourceRequirements `json:"resources" yaml:"resources"`
	// VolumeMounts: 容器的存储卷挂载
	VolumeMounts []VolumeMount `json:"volumeMounts" yaml:"volumeMounts"`

	// LivenessProbe: 存活探针，用于检测容器是否存活
	LivenessProbe *Probe `json:"livenessProbe" yaml:"livenessProbe"`
	// ReadinessProbe: 就绪探针，用于检测容器是否就绪
	ReadinessProbe *Probe `json:"readinessProbe" yaml:"readinessProbe"`
	// StartupProbe: 启动探针，用于检测容器是否启动
	StartupProbe *Probe `json:"startupProbe" yaml:"startupProbe"`
	// Lifecycle: 容器的生命周期
	Lifecycle *Lifecycle `json:"lifecycle" yaml:"lifecycle"`
	// ImagePullPolicy: 镜像拉取策略
	ImagePullPolicy string `json:"imagePullPolicy" yaml:"imagePullPolicy"`

	// Stdin: 是否使用标准输入
	Stdin bool `json:"stdin" yaml:"stdin"`
	// StdinOnce: 是否只使用一次标准输入
	StdinOnce bool `json:"stdinOnce" yaml:"stdinOnce"`
	// TTY: 是否使用 tty
	TTY bool `json:"tty" yaml:"tty"`
}

// -------------------- ContainerPort --------------------
type ContainerPort struct {
	// Name: 端口的名称
	Name string `json:"name" yaml:"name"`
	// ContainerPort: 容器端口
	ContainerPort int32 `json:"containerPort" yaml:"containerPort"`
	// Protocol: 端口协议
	Protocol Protocol `json:"protocol" yaml:"protocol"`
	// HostPort: 主机端口
	HostPort int32 `json:"hostPort" yaml:"hostPort"`
	// HostIP: 主机IP
	HostIP string `json:"hostIP" yaml:"hostIP"`
}

type Protocol string

// -------------------- EnvVar --------------------
type EnvVar struct {
	// Name: 环境变量的名称
	Name string `json:"name" yaml:"name"`
	// Value: 环境变量的值
	Value string `json:"value" yaml:"value"`
}

// -------------------- ResourceRequirements --------------------
type ResourceRequirements struct {
	// Limits: 资源限制
	Limits ResourceList `json:"limits" yaml:"limits"`
	// Requests: 资源请求
	Requests ResourceList `json:"requests" yaml:"requests"`
}

type ResourceList map[ResourceName]Quantity

// ResourceName: 资源名称，包括CPU、内存
type ResourceName string

// Quantity: 资源数量
type Quantity string

// -------------------- VolumeMount --------------------
type VolumeMount struct {
	// Name: 存储卷的名称
	Name string `json:"name" yaml:"name"`
	// MountPath: 挂载路径
	MountPath string `json:"mountPath" yaml:"mountPath"`
	// ReadOnly: 是否只读
	ReadOnly bool `json:"readOnly" yaml:"readOnly"`
}

// -------------------- Probe --------------------
type Probe struct {
	// Handler: 探针处理器
	Handler Handler `json:"handler" yaml:"handler"`
	// InitialDelaySeconds: 初始延迟时间
	InitialDelaySeconds int32 `json:"initialDelaySeconds" yaml:"initialDelaySeconds"`
	// TimeoutSeconds: 超时时间
	TimeoutSeconds int32 `json:"timeoutSeconds" yaml:"timeoutSeconds"`
	// PeriodSeconds: 周期时间
	PeriodSeconds int32 `json:"periodSeconds" yaml:"periodSeconds"`
	// SuccessThreshold: 成功阈值
	SuccessThreshold int32 `json:"successThreshold" yaml:"successThreshold"`
	// FailureThreshold: 失败阈值
	FailureThreshold int32 `json:"failureThreshold" yaml:"failureThreshold"`
}

type Handler struct {
	// Exec: 执行处理器
	Exec *ExecAction `json:"exec" yaml:"exec"`
	// HTTPGet: HTTP处理器
	HTTPGet *HTTPGetAction `json:"httpGet" yaml:"httpGet"`
	// TCPSocket: TCP处理器
	// TCPSocket *TCPSocketAction `json:"tcpSocket" yaml:"tcpSocket"`
}

type ExecAction struct {
	// Command: 执行命令
	Command []string `json:"command" yaml:"command"`
}

type HTTPGetAction struct {
	// Path: 请求路径
	Path string `json:"path" yaml:"path"`
	// Port: 请求端口
	Port int `json:"port" yaml:"port"`
	// Host: 请求主机
	Host string `json:"host" yaml:"host"`
	// Scheme: 请求协议
	Scheme string `json:"scheme" yaml:"scheme"`
}

// -------------------- Lifecycle --------------------
type Lifecycle struct {
	// PostStart: 启动后的生命周期
	PostStart *Handler `json:"postStart" yaml:"postStart"`
	// PreStop: 停止前的生命周期
	PreStop *Handler `json:"preStop" yaml:"preStop"`
}
