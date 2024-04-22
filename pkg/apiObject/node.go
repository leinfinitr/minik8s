// 描述: Node对象的封装
// 参考：https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/api/core/v1/types.go

package apiObject

type Node struct {
	// 对象的类型元数据
	TypeMeta
	// 对象的元数据
	ObjectMeta
	// Node的规格
	Spec NodeSpec
	// Node的状态
	Status NodeStatus
}

type NodeSpec struct {
	// 分配给该Node的Pod IP范围
	PodCIDR string `json:"podCIDR" yaml:"podCIDR"`
	// Node的提供者标识符
	ProviderID string `json:"providerID" yaml:"providerID"`
	// 是否可以被调度
	Unschedulable bool `json:"unschedulable" yaml:"unschedulable"`
}

type NodeStatus struct {
	// Node节点的资源容量
	Capacity map[string]string `json:"capacity" yaml:"capacity"`
	// Node的可分配资源
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
