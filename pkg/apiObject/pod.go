// 描述: Pod对象的封装
// 参考：https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/api/core/v1/types.go

package apiObject

import (
	"time"

	"github.com/docker/docker/api/types"
)

type Pod struct {
	// TypeMeta: 对象的类型元数据
	typeMeta TypeMeta
	// ObjectMeta: 对象的元数据
	objectMeta ObjectMeta
	// Spec: Pod的规格
	Spec PodSpec
	// Status: Pod的状态
	Status PodStatus
}

type PodSpec struct {
	// Volumes: Pod的存储卷
	Volumes []Volume `json:"volumes" yaml:"volumes"`
	// Containers: Pod的容器
	Containers []Container `json:"containers" yaml:"containers"`
	// RestartPolicy: 重启策略
	// 包括Always、OnFailure、Never，默认为Always
	RestartPolicy string `json:"restartPolicy" yaml:"restartPolicy"`
	// NodeSelector: 节点选择器
	// 当Pod被调度时，只有满足NodeSelector的节点才会被考虑
	NodeSelector map[string]string `json:"nodeSelector" yaml:"nodeSelector"`
	// NodeName: 表明Pod应该被调度到的节点
	// 如果为空，则表示Pod可以被调度到任何节点
	NodeName string `json:"nodeName" yaml:"nodeName"`
}

type Volume struct {
	// Name: 存储卷的名称
	Name string `json:"name" yaml:"name"`
	// VolumeSource: 存储卷的来源
	VolumeSource `json:",inline" yaml:",inline"`
}

type VolumeSource struct {
	// HostPath: 主机路径
	HostPath *HostPathVolumeSource `json:"hostPath" yaml:"hostPath"`
}

type HostPathVolumeSource struct {
	// Path: 主机路径
	Path string `json:"path" yaml:"path"`
	// Type: 主机路径类型
	Type *string `json:"type" yaml:"type"`
}

// PodStatus represents the status of a Pod.
type PodStatus struct {
	// 参考：https://kubernetes.io/zh-cn/docs/concepts/workloads/pods/pod-lifecycle/#pod-phase
	// Pending（悬决）   Pod 已被 Kubernetes 系统接受，但有一个或者多个容器尚未创建亦未运行。
	//                  此阶段包括等待 Pod 被调度的时间和通过网络下载镜像的时间。
	// Running（运行中） Pod 已经绑定到了某个节点，Pod 中所有的容器都已被创建。至少有一个容器仍在运行，
	//                   或者正处于启动或重启状态。
	// Succeeded（成功） Pod 中的所有容器都已成功终止，并且不会再重启。
	// Failed（失败）	 Pod 中的所有容器都已终止，并且至少有一个容器是因为失败终止。也就是说，容器以
	//                  非 0 状态退出或者被系统终止。
	// Unknown（未知）	 因为某些原因无法取得 Pod 的状态。这种情况通常是因为与 Pod 所在主机通信失败。
	// Terminating（需要终止） Pod 已被请求终止，但是该终止请求还没有被发送到底层容器。Pod 仍然在运行。
	Phase string `json:"conditions" yaml:"conditions"`

	// HostIP: Pod所在节点的IP地址
	HostIP string `json:"hostIP" yaml:"hostIP"`
	// HostIPs: Pod所在节点的IP地址列表
	HostIPs []string `json:"hostIPs" yaml:"hostIPs"`
	// PodIP: Pod的IP地址
	PodIP string `json:"podIP" yaml:"podIP"`
	// PodIPs: Pod的IP地址列表
	PodIPs []string `json:"podIPs" yaml:"podIPs"`

	// StartTime: Pod的启动时间
	StartTime time.Time `json:"startTime" yaml:"startTime"`
	// ContainerStatuses: 容器的状态
	ContainerStatuses []types.ContainerState `json:"containerStatuses" yaml:"containerStatuses"`
	// LastUpdateTime: 最后更新时间
	LastUpdateTime time.Time `json:"lastUpdateTime" yaml:"lastUpdateTime"`
}

func (pod *Pod) GetPodUUID() string {
	return pod.typeMeta.Metadata.UUID
}
