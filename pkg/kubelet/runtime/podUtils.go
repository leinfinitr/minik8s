package runtime

import (
	"context"
	"errors"
	"fmt"
	"minik8s/pkg/apiObject"
	"minik8s/tools/log"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

type PodUtils interface {
	CreatePod(pod *apiObject.Pod) error
	StartPod(pod *apiObject.Pod) error
	RestartPod(pod *apiObject.Pod) error
	StopPod(pod *apiObject.Pod) error
	DeletePod(pod *apiObject.Pod) error
	RecreatePodContainer(pod *apiObject.Pod) error
	ExecPodContainer(pod *apiObject.Pod) error
}

// 在这里，我们创建一个Pod相当于是创建一个Sandbox，并且会创建Pod内部的所有容器
func (r *RuntimeManager) CreatePod(pod *apiObject.Pod) error {

	sandboxConfig, err := r.getPodSandBoxConfig(pod)
	if err != nil {
		return err
	}

	request := &runtimeapi.RunPodSandboxRequest{
		Config:         sandboxConfig,
		RuntimeHandler: "",
	}

	response, err := r.runtimeClient.RunPodSandbox(context.Background(), request)
	if err != nil {
		return nil
	}

	if response.PodSandboxId == "" {
		errorMessage := fmt.Sprintf("PodSandboxId set for pod sandbox is failed ")
		log.ErrorLog(errorMessage)
		return errors.New(errorMessage)
	}
	message := fmt.Sprintf("PodSandboxId is set successfullly, Is : " + response.PodSandboxId)
	log.InfoLog(message)
	pod.PodSanboxId = response.PodSandboxId

	// 调用接口去创建Pod内部的所有容器

	return nil
}

// 创建指定配置文件的container
func (r *RuntimeManager) CreateContainers(podSandBoxID string, containerConfig *runtimeapi.ContainerConfig,
	sandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {

	request := &runtimeapi.CreateContainerRequest{
		PodSandboxId:  podSandBoxID,
		Config:        containerConfig,
		SandboxConfig: sandboxConfig,
	}

	response, err := r.runtimeClient.CreateContainer(context.Background(), request)
	if err != nil {
		errorMessage := fmt.Sprintf("create container in %d sandbox failed", podSandBoxID)
		log.ErrorLog(errorMessage)
		return "", err
	}

	return response.ContainerId, nil
}

// 运行该Pod内部的所有容器
func (r *RuntimeManager) StartPod(pod *apiObject.Pod) error {
	return nil
}

func (r *RuntimeManager) RestartPod(pod *apiObject.Pod) error {
	return nil
}

func (r *RuntimeManager) StopPod(pod *apiObject.Pod) error {
	return nil
}

func (r *RuntimeManager) DeletePod(pod *apiObject.Pod) error {
	return nil
}

func (r *RuntimeManager) RecreatePodContainer(pod *apiObject.Pod) error {
	return nil
}

func (r *RuntimeManager) ExecPodContainer(pod *apiObject.Pod) error {
	return nil
}
