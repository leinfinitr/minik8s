package apiObject

import "time"

type ScaleTargetRef struct {
	APIVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`
	MetaData ObjectMeta `json:"metadata" yaml:"metadata"`
}

type MetricSpec struct {
	CPUPercent float64 `json:"cpuPercent" yaml:"cpuPercent"`
	MemoryPercent float64 `json:"memoryPercent" yaml:"memoryPercent"`
}

type HPA struct {
	TypeMeta
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec     HPASpec    `json:"spec" yaml:"spec"`
	Status   HPAStatus  `json:"status" yaml:"status"`
}

type HPASpec struct {
	ScaletargetRef ScaleTargetRef `json:"scaleTargetRef" yaml:"scaleTargetRef"`
	MinReplicas int32 `json:"minReplicas" yaml:"minReplicas"`
	MaxReplicas int32 `json:"maxReplicas" yaml:"maxReplicas"`
	AdjustInterval time.Duration `json:"adjustInterval" yaml:"adjustInterval"`
	Metrics MetricSpec `json:"metrics" yaml:"metrics"`
}

type HPAStatus struct {
	CurrentReplicas int32 `json:"currentReplicas" yaml:"currentReplicas"`
	DesiredReplicas int32 `json:"desiredReplicas" yaml:"desiredReplicas"`
	CurCPUPercent float64 `json:"curCPUPercent" yaml:"curCPUPercent"`
	CurMemoryPercent float64 `json:"curMemoryPercent" yaml:"curMemoryPercent"`
}


