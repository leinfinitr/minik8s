package apiObject

// PersistentVolume 代表模拟的持久卷
type PersistentVolume struct {
	Name        string
	Capacity    int64 // 单位可以是字节
	Path        string
	AccessModes []string
	IsBound     bool
	ClaimedBy   string // 绑定的PVC名称
	ServerIP    string // NFS服务器的IP地址
	RemotePath  string // 在NFS服务器上的共享目录路径
}

// PersistentVolumeClaim 代表模拟的持久卷声明
type PersistentVolumeClaim struct {
	Name        string
	Namespace   string
	Capacity    int64
	AccessModes []string
	BoundTo     string // 绑定的PV名称
}
