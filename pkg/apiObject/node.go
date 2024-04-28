// 描述: Node对象的封装
// 参考：https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/api/core/v1/types.go
//		https://kubernetes.io/zh-cn/docs/reference/kubernetes-api/cluster-resources/node-v1/#NodeSpec

package apiObject

type Node struct {
	// 对象的类型元数据
	TypeMeta
	// 对象的元数据
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	// Node的规格
	Spec NodeSpec `json:"spec" yaml:"spec"`
	// Node的状态
	Status NodeStatus `json:"status" yaml:"status"`
}

type NodeSpec struct {
	// 分配给该Node的Pod IP范围
	PodCIDR string `json:"podCIDR" yaml:"podCIDR"`
	// 云提供商分配的节点ID，格式为：<ProviderName>://<ProviderSpecificNodeID>
	ProviderID string `json:"providerID" yaml:"providerID"`
	// 是否可以被调度
	Unschedulable bool `json:"unschedulable" yaml:"unschedulable"`
}

type NodeStatus struct {
	// Node节点的资源容量
	// 	包括：内存
	Capacity map[string]string `json:"capacity" yaml:"capacity"`
	// Node的可分配资源,其数值为使用率
	// 	包括：内存、CPU
	Allocatable map[string]string `json:"allocatable" yaml:"allocatable"`
	// Node最近观测到的生命周期阶段，包括：Pending、Running、Terminating
	// 	Pending: Node已经被系统创建，但是还没有被配置
	// 	Running: Node节点已经配置好，并且运行了Kubernetes组件
	// 	Terminating: Node被从集群中移除
	Phase string `json:"phase" yaml:"phase"`
	// Conditions: Node的条件
	Conditions []NodeCondition `json:"conditions" yaml:"conditions"`
	// Addresses: Node的地址
	Addresses []NodeAddress `json:"addresses" yaml:"addresses"`
}

type NodeCondition struct {
	// 条件的类型
	// 	Ready: kubelet准备好接受Pod
	// 	MemoryPressure: Node节点内存压力
	// 	DiskPressure: Node节点磁盘压力
	// 	PIDPressure: Node节点PID压力
	// 	NetworkUnavailable: Node节点网络未正确配置
	Type string `json:"type" yaml:"type"`
	// 条件的状态
	// 	True: 条件满足
	// 	False: 条件不满足
	// 	Unknown: 状态未知
	Status string `json:"status" yaml:"status"`
}

type NodeAddress struct {
	// 地址的类型
	Type string `json:"type" yaml:"type"`
	// 地址
	Address string `json:"address" yaml:"address"`
}
