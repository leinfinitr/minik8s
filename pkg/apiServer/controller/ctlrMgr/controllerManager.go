package ctlrmgr
import (
	specctlrs "minik8s/pkg/apiServer/controller/specCtlrs"
)
type ControllerManager interface {
	Run(stopCh <-chan struct{})
}
type ControllerManagerImpl struct {
	replicaSetController specctlrs.ReplicaSetController

}
func NewControllerManager() ControllerManager {
	newrc,err := specctlrs.NewReplicaController()
	if err != nil {
		panic(err)
	}
	return &ControllerManagerImpl{replicaSetController: newrc}
}

func (cm *ControllerManagerImpl) Run(stopCh <-chan struct{}) {
	go cm.replicaSetController.Run()
	<-stopCh
}
