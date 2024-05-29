package apiObject

const (
	// ReadWriteOnce 以读/写模式挂载到一个主机
	ReadWriteOnce PersistentVolumeAccessMode = "ReadWriteOnce"
	// ReadOnlyMany 以只读模式挂载到多个主机
	ReadOnlyMany PersistentVolumeAccessMode = "ReadOnlyMany"
	// ReadWriteMany 以读/写模式挂载到多个主机
	ReadWriteMany PersistentVolumeAccessMode = "ReadWriteMany"
	// ReadWriteOncePod 以读/写模式挂载到一个pod
	ReadWriteOncePod PersistentVolumeAccessMode = "ReadWriteOncePod"
)

const (
	// Retain 保留持久卷，当删除PVC的时候，PV仍然存在且不能被绑定
	Retain PersistentVolumeReclaimPolicy = "Retain"
	// Recycle 回收持久卷，当删除PVC的时候，PV会清空所有数据并准备重用
	Recycle PersistentVolumeReclaimPolicy = "Recycle"
	// Delete 删除持久卷，当删除PVC的时候，PV会被删除
	Delete PersistentVolumeReclaimPolicy = "Delete"
)

const (
	// VolumePending 用于尚未绑定且不可用的PersistentVolumes
	VolumePending PersistentVolumePhase = "Pending"
	// VolumeAvailable 用于尚未绑定且可用的PersistentVolumes
	VolumeAvailable PersistentVolumePhase = "Available"
	// VolumeBound 用于已绑定的PersistentVolumes
	VolumeBound PersistentVolumePhase = "Bound"
	// VolumeReleased 用于已绑定的PersistentVolumes，其中绑定的PersistentVolumeClaim已被删除
	//  已释放的卷必须在再次可用之前进行回收
	//  此阶段由持久卷声明绑定器用于向另一个进程发出信号以回收资源
	VolumeReleased PersistentVolumePhase = "Released"
	// VolumeFailed 用于在从声明中释放后未能正确回收或删除的PersistentVolumes
	VolumeFailed PersistentVolumePhase = "Failed"
)

type PersistentVolume struct {
	// 对象的类型元数据
	TypeMeta
	// 对象的元数据
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	// PersistentVolume的规格
	Spec PersistentVolumeSpec `json:"spec" yaml:"spec"`
	// PersistentVolume的状态
	Status PersistentVolumeStatus `json:"status" yaml:"status"`
}

type PersistentVolumeSpec struct {
	// 持久化卷的容量，单位是字节
	Capacity int64 `json:"capacity" yaml:"capacity"`
	// 访问模式
	AccessModes []PersistentVolumeAccessMode `json:"accessModes" yaml:"accessModes"`
	// 回收策略
	ReclaimPolicy PersistentVolumeReclaimPolicy `json:"reclaimPolicy" yaml:"reclaimPolicy"`
	// 远程创建持久化卷
	Remote NetworkFileSystem `json:"remote" yaml:"remote"`
}

type PersistentVolumeAccessMode string

type PersistentVolumeReclaimPolicy string

type NetworkFileSystem struct {
	// NFS服务器的IP地址
	ServerIP string `json:"serverIP" yaml:"serverIP"`
	// 在NFS服务器上的共享目录路径
	RemotePath string `json:"path" yaml:"path"`
}

type PersistentVolumeStatus struct {
	Phase PersistentVolumePhase `json:"phase" yaml:"phase"`
}

type PersistentVolumePhase string
