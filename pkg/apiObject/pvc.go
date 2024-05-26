package apiObject

// PersistentVolumeClaim 代表模拟的持久卷
type PersistentVolumeClaim struct {
	// 对象的类型元数据
	TypeMeta
	// 对象的元数据
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	// PersistentVolumeClaim的规格
	Spec PersistentVolumeClaimSpec `json:"spec" yaml:"spec"`
	// PersistentVolume的状态
	Status PersistentVolumeStatus `json:"status" yaml:"status"`
}

type PersistentVolumeClaimSpec struct {
	// 申请的持久化卷的大小，单位是字节
	Resources int64 `json:"resources" yaml:"resources"`
	// 访问模式
	AccessModes []string `json:"accessModes" yaml:"accessModes"`
}
