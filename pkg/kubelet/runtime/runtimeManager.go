package runtime

import (
	"minik8s/pkg/kubelet/runtime/dockerUtils"
	// "google.golang.org/grpc"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

type RuntimeManager struct {
	rumtimeClient    runtimeapi.RuntimeServiceClient
	containerManager dockerUtils.ContainerManager
	imageManager     dockerUtils.ImageManager
}

/* Singleton pattern */
var runtimeManager *RuntimeManager = nil

func GetRuntimeManager() *RuntimeManager {
	// TODO：从全局获取endpoint，然后需要获得与本地containerd的cnn
	if runtimeManager == nil {
		runtimeManager = &RuntimeManager{
			// runtimeManager:		runtimeapi.NewRuntimeServiceClient(nil),
			containerManager: dockerUtils.GetContainerManager(),
			imageManager:     dockerUtils.GetImageManager(),
		}
	}

	return runtimeManager
}

// func GetContainerdCnn
