package apiObject

const (
	// ClaimPending 用于尚未绑定的PersistentVolumeClaims
	ClaimPending PersistentVolumeClaimPhase = "Pending"
	// ClaimBound 用于已绑定的PersistentVolumeClaims
	ClaimBound PersistentVolumeClaimPhase = "Bound"
	// ClaimLost 用于已绑定的PersistentVolumeClaims，其中绑定的PersistentVolume已被删除
	ClaimLost PersistentVolumeClaimPhase = "Lost"
)

type PersistentVolumeClaim struct {
	// 对象的类型元数据
	TypeMeta
	// 对象的元数据
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	// PersistentVolumeClaim的规格
	Spec PersistentVolumeClaimSpec `json:"spec" yaml:"spec"`
	// PersistentVolumeClaim的状态
	Status PersistentVolumeClaimStatus `json:"status" yaml:"status"`
}

type PersistentVolumeClaimSpec struct {
	// 申请的持久化卷的大小
	Resources string `json:"resources" yaml:"resources"`
	// 访问模式
	AccessModes []PersistentVolumeAccessMode `json:"accessModes" yaml:"accessModes"`
}

type PersistentVolumeClaimStatus struct {
	// 申请的持久化卷的状态
	Phase PersistentVolumeClaimPhase `json:"phase" yaml:"phase"`
}

type PersistentVolumeClaimPhase string
