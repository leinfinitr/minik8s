// 描述: Job对象的封装
// 参考：https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/api/batch/v1/types.go

package apiObject

// Job 代表了一个Kubernetes Job对象，用于定义一次性任务（即运行至完成或失败）。
type Job struct {
	// 对象的类型元数据
	TypeMeta
	// 对象的元数据
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	// Job的规格
	Spec JobSpec `json:"spec" yaml:"spec"`
	// Job的状态
	Status JobStatus `json:"status" yaml:"status"`
}


// PodTemplateSpec 是Pod的模板定义，包含了Pod的规范和标签选择器。
type PodTemplateSpec struct {
	// 对象的元数据
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	// Pod的规格
	Spec PodSpec `json:"spec" yaml:"spec"`
}
// JobSpec 定义了Job的行为规范，如重启策略、并行度限制等。
type JobSpec struct {
	NumTasks int `json:"numTasks" yaml:"numTasks"`
	NumTasksPerNode int `json:"numTasksPerNode" yaml:"numTasksPerNode"`
	Partition string `json:"partition" yaml:"partition"`
	SubmitDir string `json:"submitDir" yaml:"submitDir"`
	UserName string `json:"userName" yaml:"userName"`
	PassWord string `json:"password" yaml:"password"`
	RunCmd []string `json:"runCmd" yaml:"runCmd"`
	OutputFile string `json:"output" yaml:"output"`
	ErrorFile string `json:"error" yaml:"error"`
	GPUNum int `json:"gpuNum" yaml:"gpuNum"`
}


// JobStatus 描述了Job当前的运行状态，如已完成Pod数量、失败次数等。
type JobStatus struct {
	JobID string `json:"jobID" yaml:"jobID"`
	Partition string `json:"partition" yaml:"partition"`
	State string `json:"state" yaml:"state"`
	ExitCode int `json:"exitCode" yaml:"exitCode"`
}

type JobCode struct{
	TypeMeta
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	UploadContent []byte `json:"uploadContent" yaml:"uploadContent"`
	OutputContent []byte `json:"outputContent" yaml:"outputContent"`
	ErrorContent []byte `json:"errorContent" yaml:"errorContent"`
}