package status

import (
	"minik8s/pkg/kubelet/runtime"
)

type StatusManager interface {
	Run()
}

type statusManagerImpl struct {
	runtimeMgr   *runtime.RuntimeManager
	apiServerURL string
}

var statusManager *statusManagerImpl = nil

func (s *statusManagerImpl) Run() {

}

// GetStatusManager 返回的是接口类型
func GetStatusManager(apiServerURL string) StatusManager {
	if statusManager == nil {
		statusManager = &statusManagerImpl{
			runtimeMgr:   runtime.GetRuntimeManager(),
			apiServerURL: apiServerURL,
		}
	}

	return statusManager
}
