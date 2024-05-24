package entity

import (
	"minik8s/pkg/apiObject"
)

const (
	CreateEvent = "CreateService"
	UpdateEvent = "UpdateService"
	DeleteEvent = "DeleteService"
)

type ServiceEvent struct {
	// 该事件的类型
	Action string
	// 该事件对应的service
	Service apiObject.Service
	// 记录service对应的所有的Endpoints
	Endpoints []apiObject.Endpoint
}
