package apiObject

const (
	ReadWriteOnce = "ReadWriteOnce"
)

// PersistentVolume 代表模拟的持久卷
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
	AccessModes []string `json:"accessModes" yaml:"accessModes"`
	// 回收策略
	PersistentVolumeReclaimPolicy string `json:"persistentVolumeReclaimPolicy" yaml:"persistentVolumeReclaimPolicy"`
	// 本地创建持久化卷
	Local LocalServer `json:"local" yaml:"local"`
	// 远程创建持久化卷
	Remote NetworkFileSystem `json:"remote" yaml:"remote"`
}

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
}
