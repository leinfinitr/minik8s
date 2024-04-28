package runtime

import (
	"context"
	"minik8s/pkg/config"
	"minik8s/pkg/klog"
	"minik8s/pkg/kubelet/util"
	"time"

	"minik8s/pkg/kubelet/runtime/dockerUtils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

type RuntimeManager struct {
	runtimeClient    runtimeapi.RuntimeServiceClient
	containerManager dockerUtils.ContainerManager
	imageManager     dockerUtils.ImageManager
}

/* Singleton pattern */
var runtimeManager *RuntimeManager = nil

func GetRuntimeManager() *RuntimeManager {
	cnn, ctx, err := GetContainerdCnn()
	if err != nil {
		return nil
	}

	// TODO：从全局获取endpoint，然后需要获得与本地containerd的cnn
	if runtimeManager == nil {
		runtimeManager = &RuntimeManager{
			runtimeClient:    runtimeapi.NewRuntimeServiceClient(cnn),
			containerManager: dockerUtils.GetContainerManager(),
			imageManager:     dockerUtils.GetImageManager(),
		}
	}

	// 我们对获取的cnn进行验证，确保能够正常使用rpc进行通信
	if _, err := runtimeManager.runtimeClient.Version(*ctx, &runtimeapi.VersionRequest{}); err != nil {
		klog.WarnLog("kubelet", "validate CRI v1 runtime API for endpoint failed")
		return nil
	}

	return runtimeManager
}

func GetContainerdCnn() (*grpc.ClientConn, *context.Context, error) {
	klog.InfoLog("kubelet", "come into GetContainerCnn")
	addr, dialer, err := util.GetAddressAndDialer(config.ContainerRuntimeEndpoint)
	if err != nil {
		return nil, nil, err
	}

	// create a context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.RuntimeRequestTimeout))
	defer cancel()

	// set some dialOpts
	var dialOpts []grpc.DialOption
	dialOpts = append(dialOpts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(config.MaxMsgSize)))

	connParams := grpc.ConnectParams{
		Backoff: backoff.DefaultConfig,
	}
	connParams.MinConnectTimeout = config.MinConnectionTimeout
	connParams.Backoff.BaseDelay = config.BaseBackoffDelay
	connParams.Backoff.MaxDelay = config.MaxBackoffDelay
	dialOpts = append(dialOpts,
		grpc.WithConnectParams(connParams),
	)

	conn, err := grpc.DialContext(ctx, addr, dialOpts...)
	if err != nil {
		klog.WarnLog("kubelet", "Connect remote runtime failed"+"address:"+addr)
		return nil, nil, err
	}

	return conn, &ctx, nil
}
