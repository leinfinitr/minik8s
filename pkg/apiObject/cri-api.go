package apiObject

type ContainerStatusRequest struct {
	// ID of the container for which to retrieve status.
	ContainerId string `protobuf:"bytes,1,opt,name=container_id,json=containerId,proto3" json:"container_id,omitempty"`
	// Verbose indicates whether to return extra information about the container.
	Verbose              bool     `protobuf:"varint,2,opt,name=verbose,proto3" json:"verbose,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}
