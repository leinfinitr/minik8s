package iptableManager

import (
	"minik8s/pkg/apiObject"
	"minik8s/pkg/entity"
)

type IptableManager interface {
	CreateService(createEvent *entity.ServiceEvent) error
	UpdateService(updateEvent *entity.ServiceEvent) error
	DeleteService(deleteEvent *entity.ServiceEvent) error
}

type iptableManager struct {

	// 实现从service的UUID到Endpoint信息的映射
	serviceUUID2Endpoint map[string][]apiObject.Endpoint
	// 从service对象到iptable内的规则链
	service2Chain map[string][]string
}

var iptableMgr *iptableManager

func GetIptableManager() IptableManager {
	if iptableMgr == nil {
		iptableMgr = &iptableManager{
			serviceUUID2Endpoint: make(map[string][]apiObject.Endpoint),
			service2Chain:        make(map[string][]string),
		}
	}

	return iptableMgr
}

func (i *iptableManager) CreateService(createEvent *entity.ServiceEvent) error {
	return nil
}

func (i *iptableManager) UpdateService(updateEvent *entity.ServiceEvent) error {
	return nil
}

func (i *iptableManager) DeleteService(deleteEvent *entity.ServiceEvent) error {
	return nil
}
