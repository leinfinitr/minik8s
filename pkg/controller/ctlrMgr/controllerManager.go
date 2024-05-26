package ctlrmgr
import (
	specctlrs "minik8s/pkg/controller/specCtlrs"
)
type ControllerManager interface {
	Run(stopCh <-chan struct{})
}
type ControllerManagerImpl struct {
	replicaSetController specctlrs.ReplicaSetController
	hpaController specctlrs.HpaController
}
func NewControllerManager() ControllerManager {
	newrc,err := specctlrs.NewReplicaController()
	if err != nil {
		panic(err)
	}
	newhc,err := specctlrs.NewHpaController()
	if err != nil {
		panic(err)
	}
	return &ControllerManagerImpl{replicaSetController: newrc, hpaController: newhc}
}

func (cm *ControllerManagerImpl) Run(stopCh <-chan struct{}) {
	go cm.replicaSetController.Run()
	<-stopCh
}
