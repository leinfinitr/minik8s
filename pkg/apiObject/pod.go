// 描述: Pod对象的封装
// 参考：https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/api/core/v1/types.go

package apiObject

import (
	"time"

	"github.com/docker/docker/api/types"
)

// 参考：https://kubernetes.io/zh-cn/docs/concepts/workloads/pods/pod-lifecycle/#pod-phase
//
//	Pending（悬决）   Pod 已被 Kubernetes 系统接受，但尚未分配至 Node 或者尚未传送至 Containerd 进行创建.
//	Running（运行中） Pod 已经绑定到了某个节点，出于创建或者重启 Pod 或者 Container 的过程中。
//	Succeeded（成功） Pod 中的所有容器都已成功运行，且正常提供 Pod 功能。
//	Failed（失败）    Pod 中的所有容器都已终止，并且至少有一个容器是因为失败终止。也就是说，容器以非 0 状态退出或者被系统终止。
//	Unknown（未知）   因为某些原因无法取得 Pod 的状态。这种情况通常是因为与 Pod 所在主机通信失败。
//	Terminating（需要终止） Pod 已被请求终止，但是该终止请求还没有被发送到底层容器。Pod 仍然在运行。
const (
	Pod_Pending     = "Pending"
	Pod_Running     = "Running"
	Pod_Successed   = "Successed"
	Pod_Failed      = "Failed"
	Pod_Unknown     = "Unknown"
	Pod_Terminating = "Terminating"
)

type Pod struct {
	// Pod对应的PodSandboxId，供查找podSandboxStatus时使用
	PodSanboxId string
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
	// 存储卷的来源
	VolumeSource `json:",inline" yaml:",inline"`
}

type VolumeSource struct {
	// 主机路径
	HostPath *HostPathVolumeSource `json:"hostPath" yaml:"hostPath"`
}

type HostPathVolumeSource struct {
	// 主机路径
	Path string `json:"path" yaml:"path"`
	// 主机路径类型
	Type *string `json:"type" yaml:"type"`
}

// PodStatus represents the status of a Pod.
type PodStatus struct {
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
}

func (p *Pod) GetPodUUID() string {
	return p.Metadata.UUID
}
