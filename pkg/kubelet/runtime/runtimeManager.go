package runtime

import (
	"minik8s/pkg/config"
	"minik8s/tools/log"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

type RuntimeManager struct {
	runtimeClient runtimeapi.RuntimeServiceClient
	imageManager  ImageManager
}

/* Singleton pattern */
var runtimeManager *RuntimeManager = nil

func GetRuntimeManager() *RuntimeManager {
	cnn, ctx, cancel, err := GetCnn(config.ContainerRuntimeEndpoint)
	if err != nil {
		return nil
	}
	defer cancel()

	if runtimeManager == nil {
		runtimeManager = &RuntimeManager{
			runtimeClient: runtimeapi.NewRuntimeServiceClient(cnn),
			imageManager:  GetImageManager(),
		}
	}

	// 我们对获取的cnn进行验证，确保能够正常使用rpc进行通信
	if _, err := runtimeManager.runtimeClient.Version(*ctx, &runtimeapi.VersionRequest{}); err != nil {
		log.ErrorLog("validate CRI v1 runtime API for endpoint failed")
		return nil
	}

	return runtimeManager
}
