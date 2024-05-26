// 描述: Pod对象的封装
// 参考：https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/api/core/v1/types.go

package apiObject

import (
	"time"

	"github.com/docker/docker/api/types"
)

// 参考：https://kubernetes.io/zh-cn/docs/concepts/workloads/pods/pod-lifecycle/#pod-phase
//
//		Pending（悬决）   Pod 已被 Kubernetes 系统接受，但尚未分配至 Node 或者尚未传送至 Containerd 进行创建.
//	 	Building (创建中) Pod 已经被分配到具体的 Node 节点进行创建，此时正在创建 PodSandbox 或者 Containers，或者正在准备镜像
//		Succeeded（成功创建） Pod 中的所有容器都已成功创建，但是Pod还未被运行。
//		Running（运行中） Pod 已经成功运行，且正常提供 Pod 功能，有kubelet负责其容错。
//		Failed（失败）    Pod 中的所有容器都已终止，并且至少有一个容器是因为失败终止。也就是说，容器以非 0 状态退出或者被系统终止。
//		Unknown（未知）   因为某些原因无法取得 Pod 的状态。这种情况通常是因为与 Pod 所在主机通信失败。
//		Terminating（需要终止） Pod 已被请求终止，但是该终止请求还没有被发送到底层容器。Pod 仍然在运行。
const (
	PodPending     = "Pending"
	PodBuilding    = "Building"
	PodRunning     = "Running"
	PodSucceeded   = "Succeeded"
	PodFailed      = "Failed"
	PodUnknown     = "Unknown"
	PodTerminating = "Terminating"
)

type Pod struct {
	// Pod对应的PodSandboxId，供查找podSandboxStatus时使用
	PodSandboxId string
	// 对象的类型元数据
	TypeMeta
	// 对象的元数据
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	// Pod的规格
	Spec PodSpec `json:"spec" yaml:"spec"`
	// Pod的状态
	Status PodStatus `json:"status" yaml:"status"`
}

type PodSpec struct {
	// Pod的存储卷
	Volumes []Volume `json:"volumes" yaml:"volumes"`
	// Pod的容器
	Containers []Container `json:"containers" yaml:"containers"`
	// 重启策略, 包括Always、OnFailure、Never，默认为Always
	RestartPolicy string `json:"restartPolicy" yaml:"restartPolicy"`
	// 节点选择器
	// 	当Pod被调度时，只有满足NodeSelector的节点才会被考虑
	NodeSelector map[string]string `json:"nodeSelector" yaml:"nodeSelector"`
	// 表明Pod应该被调度到的节点
	// 	如果为空，则表示Pod可以被调度到任何节点
	NodeName string `json:"nodeName" yaml:"nodeName"`
}

type Volume struct {
	// 存储卷的名称
	Name string `json:"name" yaml:"name"`
	// 空目录
	EmptyDir EmptyDirVolumeSource `json:"emptyDir" yaml:"emptyDir"`
	// 存储卷的来源
	VolumeSource
	// 持久化卷声明
	PersistentVolumeClaim PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim" yaml:"persistentVolumeClaim"`
}

type EmptyDirVolumeSource struct {
	Medium string `json:"medium" yaml:"medium"`
	// 空目录的大小
	SizeLimit string `json:"sizeLimit" yaml:"sizeLimit"`
}
type VolumeSource struct {
	// 主机路径
	HostPath HostPathVolumeSource `json:"hostPath" yaml:"hostPath"`
}

type HostPathVolumeSource struct {
	// 主机路径
	Path string `json:"path" yaml:"path"`
	// 主机路径类型
	Type string `json:"type" yaml:"type"`
}

type PersistentVolumeClaimVolumeSource struct {
	// 持久化卷声明的名称
	ClaimName string `json:"claimName" yaml:"claimName"`
	// 持久化卷声明的读写模式
	ReadOnly bool `json:"readOnly" yaml:"readOnly"`
}

// PodStatus represents the status of a Pod.
type PodStatus struct {
	// Pod当前所处的阶段
	Phase string `json:"conditions" yaml:"conditions"`

	// Pod所在节点的IP地址
	HostIP string `json:"hostIP" yaml:"hostIP"`
	// Pod所在节点的IP地址列表
	HostIPs []string `json:"hostIPs" yaml:"hostIPs"`
	// Pod的IP地址
	PodIP string `json:"podIP" yaml:"podIP"`
	// Pod的IP地址列表
	PodIPs []string `json:"podIPs" yaml:"podIPs"`

	// Pod的启动时间
	StartTime time.Time `json:"startTime" yaml:"startTime"`
	// 容器的状态
	ContainerStatuses []types.ContainerState `json:"containerStatuses" yaml:"containerStatuses"`
	// 最后更新时间
	LastUpdateTime time.Time `json:"lastUpdateTime" yaml:"lastUpdateTime"`

	//Pod资源状态
	CpuUsage float64 `json:"cpuUsage" yaml:"cpuUsage"`
	MemUsage float64 `json:"memUsage" yaml:"memUsage"`
}

func (p *Pod) GetPodUUID() string {
	return p.Metadata.UUID
}
