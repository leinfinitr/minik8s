package dockerUtils

import (
	"minik8s/pkg/apiObject"
	// dockerClient "github.com/docker/docker/client"
)

/* containerManager 负责 node 节点上运行的容器的 cgroup 配置信息，
 * kubelet 启动参数如果指定 --cgroups-per-qos 的时候，kubelet 会启动
 * goroutine 来周期性的更新 pod 的 cgroup 信息，维护其正确性，该参数默认为 true,
 * 实现了 pod 的Guaranteed/BestEffort/Burstable 三种级别的 Qos。
 */

type ContainerManager interface {
	CreateContainer(name string, options *apiObject.Container) (string, error)
}

type containerManagerImpl struct {
}

var containerManager *containerManagerImpl = nil

func GetContainerManager() ContainerManager {
	if containerManager == nil {
		containerManager = &containerManagerImpl{}
	}
	return containerManager
}

func (c *containerManagerImpl) CreateContainer(name string, options *apiObject.Container) (string, error) {
	return "", nil
	/* Step-1: get a new docker client */
	// if client, err := NewDockerClient; err != nil {
	// return "", err
	// }

	/* Step-2: pull the image if needed */
	// imageManager := &dockerUtils.ImageManager

	/* Step-3: */
}

/* 创建一个 docker client 对象 */
// func NewDockerClient() (*dockerClient.Client, error) {
// 	cli, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithAPIVersionNegotiation())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return cli, nil
// }
