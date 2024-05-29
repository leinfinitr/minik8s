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
	PersistentVolumeReclaimPolicy string `json:"persistentVolumeReclaimPolicy" yaml:"persistentVolumeReclaimPolicy"`
	// 本地创建持久化卷
	Local LocalServer `json:"local" yaml:"local"`
	// 远程创建持久化卷
	Remote NetworkFileSystem `json:"remote" yaml:"remote"`
}

type PersistentVolumeAccessMode string

type LocalServer struct {
	// 本机存储路径
	Path string `json:"path" yaml:"path"`
}

type NetworkFileSystem struct {
	// NFS服务器的IP地址
	ServerIP string `json:"serverIP" yaml:"serverIP"`
	// 在NFS服务器上的共享目录路径
	RemotePath string `json:"remotePath" yaml:"remotePath"`
}

type PersistentVolumeStatus struct {
	Phase PersistentVolumePhase `json:"phase" yaml:"phase"`
}

type PersistentVolumePhase string
