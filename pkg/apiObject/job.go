// 描述: Job对象的封装
// 参考：https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/api/batch/v1/types.go

package apiObject

// Job 代表了一个Kubernetes Job对象，用于定义一次性任务（即运行至完成或失败）。
type Job struct {
	// 对象的类型元数据
	TypeMeta
	// 对象的元数据
	ObjectMeta
	// Job的规格
	Spec JobSpec
	// Job的状态
	Status JobStatus
}

// JobSpec 定义了Job的行为规范，如重启策略、并行度限制等。
type JobSpec struct {
	// 控制Job并行执行Pod的最大数量。
	Parallelism *int32 `json:"parallelism" yaml:"parallelism"`
	// Job期望完成的Pod副本数量。当这个数量的Pod成功完成时，Job被认为完成。
	Completions int32 `json:"completions" yaml:"completions"`
	// 超时时间，在此时间内Job必须完成，否则会被视为失败。
	ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds" yaml:"activeDeadlineSeconds"`
	// 指定Pod模板，Job会基于此模板创建Pods来执行任务。
	Template PodTemplateSpec `json:"template" yaml:"template"`
}

// PodTemplateSpec 是Pod的模板定义，包含了Pod的规范和标签选择器。
type PodTemplateSpec struct {
	// 对象的元数据
	ObjectMeta
	// Pod的规格
	Spec PodSpec
}

// JobStatus 描述了Job当前的运行状态，如已完成Pod数量、失败次数等。
type JobStatus struct {
	// 完成的Pod总数。
	Succeeded int32 `json:"succeeded,omitempty" yaml:"succeeded,omitempty"`
	// 失败的Pod总数。
	Failed int32 `json:"failed,omitempty" yaml:"failed,omitempty"`
}
