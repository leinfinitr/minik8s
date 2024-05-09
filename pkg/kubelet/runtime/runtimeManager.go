package runtime

import (
	"minik8s/pkg/config"
	"minik8s/tools/log"

	"minik8s/pkg/kubelet/runtime/dockerUtils"
	"minik8s/pkg/kubelet/runtime/image"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

type RuntimeManager struct {
	runtimeClient    runtimeapi.RuntimeServiceClient
	containerManager dockerUtils.ContainerManager
	imageManager     image.ImageManager
}

/* Singleton pattern */
var runtimeManager *RuntimeManager = nil

func GetRuntimeManager() *RuntimeManager {
	cnn, ctx, cancel, err := image.GetCnn(config.ContainerRuntimeEndpoint)
	if err != nil {
		return nil
	}
	defer cancel()

	// TODO：从全局获取endpoint，然后需要获得与本地containerd的cnn
	if runtimeManager == nil {
		runtimeManager = &RuntimeManager{
			runtimeClient:    runtimeapi.NewRuntimeServiceClient(cnn),
			containerManager: dockerUtils.GetContainerManager(),
			imageManager:     image.GetImageManager(),
		}
	}

	// 我们对获取的cnn进行验证，确保能够正常使用rpc进行通信
	if _, err := runtimeManager.runtimeClient.Version(*ctx, &runtimeapi.VersionRequest{}); err != nil {
		log.WarnLog("validate CRI v1 runtime API for endpoint failed")
		return nil
	}

	return runtimeManager
}
