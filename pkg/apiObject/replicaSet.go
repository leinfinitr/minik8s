// 描述: replicaSet对象的定义
// 参考：https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/api/apps/v1/types.go

package apiObject

type ReplicaSet struct {
	// 对象的类型元数据
	TypeMeta
	// 对象的元数据
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	// ReplicaSet的规格
	Spec ReplicaSetSpec `json:"spec" yaml:"spec"`
	// ReplicaSet的状态
	Status ReplicaSetStatus `json:"status" yaml:"status"`
}

type ReplicaSetSpec struct {
	// 所需replica的数量
	Replicas int32 `json:"replicas" yaml:"replicas"`
	// 新创建的pod在没有任何容器崩溃的情况下准备就绪的最小秒数，以便被视为可用
	MinReadySeconds int32 `json:"minReadySeconds" yaml:"minReadySeconds"`
	// 选择器，用于标识ReplicaSet管理的Pod
	Selector map[string]string `json:"selector" yaml:"selector"`
	// 在没有足够的replica时，pod的创建模板
	Template PodTemplateSpec `json:"template" yaml:"template"`
}

type ReplicaSetStatus struct {
	// 已经创建的replica的数量
	Replicas int32 `json:"replicas" yaml:"replicas"`
	// 已经就绪的replica的数量
	ReadyReplicas int32 `json:"readyReplicas" yaml:"readyReplicas"`
	// 副本集当前状态的最新可用观测值
	Conditions []ReplicaSetCondition `json:"conditions" yaml:"conditions"`
}

type ReplicaSetCondition struct {
	// 条件的类型
	// 	ReplicaFailure: 副本失败
	// 	ReplicaSuccess: 副本成功
	Type string `json:"type" yaml:"type"`
	// 条件的状态
	// 	True: 条件满足
	// 	False: 条件不满足
	// 	Unknown: 状态未知
	Status string `json:"status" yaml:"status"`
	// 最后一次状态变化的时间
	LastTransitionTime string `json:"lastTransitionTime" yaml:"lastTransitionTime"`
	// 条件最后一次转换的原因
	Reason string `json:"reason" yaml:"reason"`
	// 条件的详细信息
	Message string `json:"message" yaml:"message"`
}
