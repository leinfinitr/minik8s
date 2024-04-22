// 描述: Job对象的封装
// 参考：https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/api/batch/v1/types.go

package apiObject

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

type JobSpec struct {
	// TODO：根据实际的需求定义Job的规格
}

type JobStatus struct {
	// TODO：根据实际的需求定义Job的状态
}
