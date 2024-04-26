package runtime

import (
	"minik8s/pkg/kubelet/runtime/dockerUtils"
)

type RuntimeManager struct {
	containerManager dockerUtils.ContainerManager
	imageManager     dockerUtils.ImageManager
}

/* Singleton pattern */
var runtimeManager *RuntimeManager = nil

func GetRuntimeManager() *RuntimeManager {
	if runtimeManager == nil {
		runtimeManager = &RuntimeManager{
			containerManager: dockerUtils.GetContainerManager(),
			imageManager:     dockerUtils.GetImageManager(),
		}
	}

	return runtimeManager
}
