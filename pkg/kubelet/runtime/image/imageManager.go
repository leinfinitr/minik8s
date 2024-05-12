package image

import (
	"context"
	"strings"
	"time"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/pkg/kubelet/util"
	"minik8s/tools/log"

	dockerref "github.com/distribution/reference"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

/* imageManager 调用 借助docker来实现对于容器镜像的管理
 */
type ImageManager interface {
	PullImage(container *apiObject.Container, sandboxConfig *runtimeapi.PodSandboxConfig) (string, error)
	ImageStatus(image *runtimeapi.ImageSpec, ifVerbose bool) (string, error)
}

type imageManagerImpl struct {
	ImageClient runtimeapi.ImageServiceClient
}

var imageManager *imageManagerImpl = nil

func GetImageManager() ImageManager {
	cnn, ctx, cancel, err := GetCnn(config.ImageRuntimeEndpoint)
	if err != nil {
		return nil
	}
	defer cancel()

	if imageManager == nil {
		imageManager = &imageManagerImpl{
			ImageClient: runtimeapi.NewImageServiceClient(cnn),
		}
	}

	// 对获得的client进行通信测试
	if _, err := imageManager.ImageClient.ImageFsInfo(*ctx, &runtimeapi.ImageFsInfoRequest{}); err != nil {
		log.ErrorLog("Image client cnn test failed")
		return nil
	}

	return imageManager
}

// 查找镜像当前的状态，如果镜像已经在本地，则返回镜像的ID，否则返回空字符串
func (i *imageManagerImpl) ImageStatus(image *runtimeapi.ImageSpec, ifVerbose bool) (string, error) {
	request := &runtimeapi.ImageStatusRequest{
		Image:   image,
		Verbose: ifVerbose,
	}

	response, err := i.ImageClient.ImageStatus(context.Background(), request)

	if err != nil {
		log.ErrorLog("[RPC] get image status failed")
		return "", err
	}

	if response.Image == nil || response.Image.Id == "" || response.Image.Size_ == 0 {
		return "", nil
	}

	return response.Image.Id, nil
}

// 返回镜像的ImageRef(也就会ImageID)，如果有需要则需要从仓库中拉下镜像
func (i *imageManagerImpl) PullImage(container *apiObject.Container, sandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {

	image, err := generateDefaultImageTag(container.Image)
	if err != nil {
		log.ErrorLog("generate imageRef failed")
		return "", err
	}

	var annotations map[string]string

	spec := &runtimeapi.ImageSpec{
		Image:          image,
		Annotations:    annotations,
		RuntimeHandler: "", // TODO: 这里的RuntimeHandler是否会用的到？？
	}

	// 先查找该镜像是否已经被在本地了
	imageID, err := i.ImageStatus(spec, false)
	if err != nil {
		return "", nil
	}
	if imageID != "" {
		// 说明镜像已经在本地了，返回imageID
		return imageID, nil
	}

	// 不在本地，则需要使用client从云端仓库拉取到本地(此处我们忽略下拉策略的具体配置)

	request := &runtimeapi.PullImageRequest{
		Image:         spec,
		Auth:          nil,
		SandboxConfig: sandboxConfig,
	}
	response, err := i.ImageClient.PullImage(context.Background(), request)
	if err != nil {
		log.ErrorLog("[RPC] Pull Image failed: " + err.Error())
		return "", nil
	}

	return response.ImageRef, nil
}

func generateDefaultImageTag(image string) (string, error) {
	named, err := dockerref.ParseNormalizedNamed(image)
	if err != nil {
		return "", nil
	}

	var tag, digest string

	tagged, ok := named.(dockerref.Tagged)
	if ok {
		tag = tagged.Tag()
	}

	digested, ok := named.(dockerref.Digested)
	if ok {
		digest = digested.Digest().String()
	}
	// If no tag was specified, use the default "latest".
	if len(tag) == 0 && len(digest) == 0 {
		tag = "latest"
	}

	if len(digest) == 0 && len(tag) > 0 && !strings.HasSuffix(image, ":"+tag) {
		image = image + ":" + tag
	}
	return image, nil

}

func GetCnn(endpoint string) (*grpc.ClientConn, *context.Context, context.CancelFunc, error) {
	log.InfoLog("come into GetContainerCnn")
	addr, dialer, err := util.GetAddressAndDialer(endpoint)
	if err != nil {
		return nil, nil, nil, err
	}

	// create a context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.RuntimeRequestTimeout))

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
		log.WarnLog("Connect remote runtime failed" + "address:" + addr)
		return nil, nil, nil, err
	}

	return conn, &ctx, cancel, nil
}
