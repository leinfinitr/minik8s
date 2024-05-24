package apiObject

type ExecReq struct {
	// ID of the container in which to execute the command.
	ContainerId string `protobuf:"bytes,1,opt,name=container_id,json=containerId,proto3" json:"container_id,omitempty"`
	// Command to execute.
	Cmd []string `protobuf:"bytes,2,rep,name=cmd,proto3" json:"cmd,omitempty"`
	// Whether to exec the command in a TTY.
	Tty bool `protobuf:"varint,3,opt,name=tty,proto3" json:"tty,omitempty"`
	// Whether to stream stdin.
	// One of `stdin`, `stdout`, and `stderr` MUST be true.
	Stdin bool `protobuf:"varint,4,opt,name=stdin,proto3" json:"stdin,omitempty"`
	// Whether to stream stdout.
	// One of `stdin`, `stdout`, and `stderr` MUST be true.
	Stdout bool `protobuf:"varint,5,opt,name=stdout,proto3" json:"stdout,omitempty"`
	// Whether to stream stderr.
	// One of `stdin`, `stdout`, and `stderr` MUST be true.
	// If `tty` is true, `stderr` MUST be false. Multiplexing is not supported
	// in this case. The output of stdout and stderr will be combined to a
	// single stream.
	Stderr bool `protobuf:"varint,6,opt,name=stderr,proto3" json:"stderr,omitempty"`
}

type ExecRsp struct {
	// Fully qualified URL of the exec streaming server.
	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}
