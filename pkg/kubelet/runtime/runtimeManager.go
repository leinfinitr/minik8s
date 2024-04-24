package runtime

import (
	"minik8s/pkg/kubelet/runtime/container"
	"minik8s/pkg/kubelet/runtime/image"
)

type RuntimeManager struct {
	containerManager container.ContainerManager
	imageManager     image.ImageManager
}

/* Singleton pattern */
var runtimeManager *RuntimeManager = nil

func GetRuntimeManager() *RuntimeManager {
	if runtimeManager == nil {
		runtimeManager = &RuntimeManager{
			containerManager: container.ContainerManager{},
			imageManager:     image.ImageManager{},
		}
	}

	return runtimeManager
}
