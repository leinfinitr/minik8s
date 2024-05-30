// 描述: 容器相关的 API 对象
// 参考：https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/api/core/v1/types.go

package apiObject

type ContainerStatus int

const (
	ContainerUncreated ContainerStatus = -1
	ContainerCreated   ContainerStatus = 0
	ContainerRunning   ContainerStatus = 1
	ContainerExited    ContainerStatus = 2
	ContainerUnknown   ContainerStatus = 3
)

const (
	// CPU, 单位 cores. (500m = .5 cores)
	ResourceCPU ResourceName = "cpu"
	// Memory, 单位 bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceMemory ResourceName = "memory"
	// storage, 单位 bytes (e,g. 5Gi = 5GiB = 5 * 1024 * 1024 * 1024)
	ResourceStorage ResourceName = "storage"
	// 本地临时存储, 单位 bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceEphemeralStorage ResourceName = "ephemeral-storage"
)

type Container struct {
	// 容器的ID
	ContainerID string
	// 容器当前状态
	ContainerStatus ContainerStatus
	// 容器的名称
	Name string `json:"name" yaml:"name"`
	// 容器的镜像
	Image string `json:"image" yaml:"image"`
	// 容器的命令
	Command []string `json:"command" yaml:"command"`
	// 容器的命令行参数
	Args []string `json:"args" yaml:"args"`

	// 容器的工作目录
	WorkingDir string `json:"workingDir" yaml:"workingDir"`
	// 容器的端口
	Ports []ContainerPort `json:"ports" yaml:"ports"`
	// 容器的环境变量
	Env []EnvVar `json:"env" yaml:"env"`
	// 容器的资源限制
	Resources ResourceRequirements `json:"resources" yaml:"resources"`
	// 容器的存储卷挂载
	VolumeMounts []VolumeMount `json:"volumeMounts" yaml:"volumeMounts"`
	// 容器与主机的挂载
	Mounts []*Mount `json:"mounts" yaml:"mounts"`

	// 镜像拉取策略
	ImagePullPolicy string `json:"imagePullPolicy" yaml:"imagePullPolicy"`

	// 是否使用标准输入
	Stdin bool `json:"stdin" yaml:"stdin"`
	// 是否只使用一次标准输入
	StdinOnce bool `json:"stdinOnce" yaml:"stdinOnce"`
	// 是否使用 tty
	TTY bool `json:"tty" yaml:"tty"`
}

type ContainerPort struct {
	// 端口的名称
	Name string `json:"name" yaml:"name"`
	// 是否是metrics端口，如果不为空，则表示是metrics端口，会进行监控和管理
	Metrics string
	// 容器端口
	ContainerPort int32 `json:"containerPort" yaml:"containerPort"`
	// 端口协议
	Protocol Protocol `json:"protocol" yaml:"protocol"`
	// 主机端口
	HostPort int32 `json:"hostPort" yaml:"hostPort"`
	// 主机IP
	HostIP string `json:"hostIP" yaml:"hostIP"`
}

// 端口协议
type Protocol string

type EnvVar struct {
	// 环境变量的名称
	Name string `json:"name" yaml:"name"`
	// 环境变量的值
	Value string `json:"value" yaml:"value"`
}

type ResourceRequirements struct {
	// 资源限制
	Limits ResourceList `json:"limits" yaml:"limits"`
	// 资源请求
	Requests ResourceList `json:"requests" yaml:"requests"`
}

// 资源列表
type ResourceList map[ResourceName]string

// 资源名称，包括CPU、内存
type ResourceName string

type VolumeMount struct {
	// 存储卷的名称
	Name string `json:"name" yaml:"name"`
	// 挂载路径
	MountPath string `json:"mountPath" yaml:"mountPath"`
	// 是否只读
	ReadOnly bool `json:"readOnly" yaml:"readOnly"`
}

type Mount struct {
	// 主机路径
	HostPath string `json:"hostPath" yaml:"hostPath"`
	// 容器路径
	ContainerPath string `json:"containerPath" yaml:"containerPath"`
	// 是否只读
	ReadOnly bool `json:"readOnly" yaml:"readOnly"`
}

type ExecAction struct {
	// 执行命令
	Command []string `json:"command" yaml:"command"`
}

type HTTPGetAction struct {
	// 请求路径
	Path string `json:"path" yaml:"path"`
	// 请求端口
	Port int `json:"port" yaml:"port"`
	// 请求主机
	Host string `json:"host" yaml:"host"`
	// 请求协议
	Scheme string `json:"scheme" yaml:"scheme"`
}
