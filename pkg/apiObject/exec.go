package apiObject

type ExecReq struct {
	// 执行命令的容器ID
	ContainerId string `protobuf:"bytes,1,opt,name=container_id,json=containerId,proto3" json:"container_id,omitempty"`
	// 命令
	Cmd []string `protobuf:"bytes,2,rep,name=cmd,proto3" json:"cmd,omitempty"`
	// 是否使用tty
	Tty bool `protobuf:"varint,3,opt,name=tty,proto3" json:"tty,omitempty"`
	// 是否使用标准输入
	Stdin bool `protobuf:"varint,4,opt,name=stdin,proto3" json:"stdin,omitempty"`
	// 是否使用标准输出
	Stdout bool `protobuf:"varint,5,opt,name=stdout,proto3" json:"stdout,omitempty"`
	// 是否使用标准错误输出，如果使用tty，stderr必须为false
	Stderr bool `protobuf:"varint,6,opt,name=stderr,proto3" json:"stderr,omitempty"`
}

type ExecRsp struct {
	// 执行流式服务器的完全限定URL
	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}
