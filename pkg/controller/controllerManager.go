package controller

import (
	specctlrs "minik8s/pkg/controller/specCtlrs"
	"minik8s/tools/log"
)

type ControllerManager interface {
	Run(stopCh <-chan struct{})
}

type ControllerManagerImpl struct {
	replicaSetController specctlrs.ReplicaSetController
	hpaController        specctlrs.HpaController
	pvController         specctlrs.PvController
	jobController 	  specctlrs.JobController
}

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
	newjc, err := specctlrs.NewJobController()
	if err != nil {
		panic(err)
	}
	return &ControllerManagerImpl{replicaSetController: newrc, hpaController: newhc, pvController: newpc, jobController: newjc}
}

func (cm *ControllerManagerImpl) Run(stopCh <-chan struct{}) {
	log.InfoLog("ControllerManager Run")
	go cm.replicaSetController.Run()
	go cm.hpaController.Run()
	go cm.pvController.Run()
	go cm.jobController.Run()
	<-stopCh
}
