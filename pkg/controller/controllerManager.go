package controller

import (
	"minik8s/pkg/apiObject"
	specctlrs "minik8s/pkg/controller/specCtlrs"
)

type ControllerManager interface {
	Run(stopCh <-chan struct{})

	AddPv(pv *apiObject.PersistentVolume) error
	AddPvc(pvc *apiObject.PersistentVolumeClaim) error
	BindPodToPvc(pvc *apiObject.PersistentVolumeClaim, podName string) error
	UnbindPodToPvc(pvc *apiObject.PersistentVolumeClaim) error
	GetPvcBind(pvcName string) string
}

type ControllerManagerImpl struct {
	replicaSetController specctlrs.ReplicaSetController
	hpaController        specctlrs.HpaController
	pvController         specctlrs.PvController
}

var ControllerManagerInstance *ControllerManagerImpl = nil

func NewControllerManager() ControllerManager {
	newrc, err := specctlrs.NewReplicaController()
	if err != nil {
		panic(err)
	}
	newhc, err := specctlrs.NewHpaController()
	if err != nil {
		panic(err)
	}
	newpc, err := specctlrs.NewPvController()
	if err != nil {
		panic(err)
	}
	if ControllerManagerInstance == nil {
		ControllerManagerInstance = &ControllerManagerImpl{replicaSetController: newrc, hpaController: newhc, pvController: newpc}
	}
	return ControllerManagerInstance
}

func (cm *ControllerManagerImpl) Run(stopCh <-chan struct{}) {
	go cm.replicaSetController.Run()
	go cm.hpaController.Run()
	go cm.pvController.Run()
	<-stopCh
}

func (cm *ControllerManagerImpl) AddPv(pv *apiObject.PersistentVolume) error {
	return cm.pvController.AddPv(pv)
}

func (cm *ControllerManagerImpl) AddPvc(pvc *apiObject.PersistentVolumeClaim) error {
	return cm.pvController.AddPvc(pvc)
}

func (cm *ControllerManagerImpl) GetPvcBind(pvcName string) string {
	return cm.pvController.GetPvcBind(pvcName)
}

func (cm *ControllerManagerImpl) BindPodToPvc(pvc *apiObject.PersistentVolumeClaim, podName string) error {
	return cm.pvController.BindPodToPvc(pvc, podName)
}

func (cm *ControllerManagerImpl) UnbindPodToPvc(pvc *apiObject.PersistentVolumeClaim) error {
	return cm.pvController.UnbindPodToPvc(pvc)
}
