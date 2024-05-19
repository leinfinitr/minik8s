package specctlrs

type ReplicaSetController interface {
	Run()
}
type ReplicaSetControllerImpl struct {
}

func NewReplicaController() (ReplicaSetController, error) {
	return &ReplicaSetControllerImpl{}, nil
}

func (rc *ReplicaSetControllerImpl) Run() {
	// 定期执行
	//TODO
}