package runtime

import (
	"context"
	"errors"
	"fmt"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
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
	log.InfoLog("[RPC] Start CreatePod")

	sandboxConfig, err := r.getPodSandBoxConfig(pod)
	if err != nil {
		log.ErrorLog("GetPodSandBoxConfig fail: " + err.Error())
		return err
	}

	request := &runtimeapi.RunPodSandboxRequest{
		Config:         sandboxConfig,
		RuntimeHandler: "",
	}

	log.DebugLog("RunPodSandbox")
	log.DebugLog(request.String())
	response, err := r.runtimeClient.RunPodSandbox(context.Background(), request)
	if err != nil {
		log.ErrorLog("Create Podsandbox fail " + err.Error())
		return nil
	}

	log.DebugLog("Create Podsandbox Success")

	if response.PodSandboxId == "" {
		errorMessage := "PodSandboxId set for pod sandbox is failed "
		log.ErrorLog(errorMessage)
		return errors.New(errorMessage)
	}
	message := fmt.Sprintf("PodSandboxId is set successfullly, Is : " + response.PodSandboxId)
	log.InfoLog(message)
	pod.PodSanboxId = response.PodSandboxId

	// 调用接口去创建Pod内部的所有容器
	containers := &pod.Spec.Containers
	for i := 0; i < len(*containers); i += 1 {
		containerConfig, err := r.getContainerConfig(&(*containers)[i], sandboxConfig)
		if err != nil {
			log.ErrorLog("generate container config failed")
			return err
		}

		containerID, err := r.CreateContainers(pod.PodSanboxId, containerConfig, sandboxConfig)
		if err != nil {
			log.ErrorLog("Create containers failed")
			return err
		}

		(*containers)[i].ContainerID = containerID
		(*containers)[i].ContainerStatus = apiObject.Container_Created
		message := fmt.Sprintf("container Id:%s is created ", containerID)
		log.InfoLog(message)
	}

	return nil
}

// 创建指定配置文件的container
func (r *RuntimeManager) CreateContainers(podSandBoxID string, containerConfig *runtimeapi.ContainerConfig,
	sandboxConfig *runtimeapi.PodSandboxConfig) (string, error) {
	log.InfoLog("[RPC] Start CreateContainers")
	request := &runtimeapi.CreateContainerRequest{
		PodSandboxId:  podSandBoxID,
		Config:        containerConfig,
		SandboxConfig: sandboxConfig,
	}

	response, err := r.runtimeClient.CreateContainer(context.Background(), request)
	if err != nil {
		errorMessage := fmt.Sprintf("create container in %s sandbox failed", podSandBoxID)
		log.ErrorLog(errorMessage)
		return "", err
	}

	return response.ContainerId, nil
}

// 运行该Pod内部的所有容器
func (r *RuntimeManager) StartPod(pod *apiObject.Pod) error {
	log.InfoLog("[RPC] Start StartPod")
	for i := 0; i < len(pod.Spec.Containers); i += 1 {
		_, err := r.runtimeClient.StartContainer(context.Background(), &runtimeapi.StartContainerRequest{
			ContainerId: pod.Spec.Containers[i].ContainerID,
		})
		if err != nil {
			errorMsg := fmt.Sprintf("[RPC] Start container failed, containerID: %s", pod.Spec.Containers[i].ContainerID)
			log.ErrorLog(errorMsg)
			return err
		}
		pod.Spec.Containers[i].ContainerStatus = apiObject.Container_Running

	}
	return nil
}

func (r *RuntimeManager) RestartPod(pod *apiObject.Pod) error {
	log.InfoLog("[RPC] Start RestartPod")
	// 考虑到容器之间可能存在依赖，为了保证可用性，在暂停所有的容器后再重新启动
	for i := 0; i < len(pod.Spec.Containers); i += 1 {
		_, err := r.runtimeClient.StopContainer(context.Background(), &runtimeapi.StopContainerRequest{
			ContainerId: pod.Spec.Containers[i].ContainerID,
			Timeout:     config.RPCRequestTimeout,
		})
		if err != nil {
			errorMsg := fmt.Sprintf("[RPC] In restarting, stop container failed, containerID: %s", pod.Spec.Containers[i].ContainerID)
			log.ErrorLog(errorMsg)
			return err
		}

		pod.Spec.Containers[i].ContainerStatus = apiObject.Container_Restart
	}

	for i := 0; i < len(pod.Spec.Containers); i += 1 {
		_, err := r.runtimeClient.StartContainer(context.Background(), &runtimeapi.StartContainerRequest{
			ContainerId: pod.Spec.Containers[i].ContainerID,
		})
		if err != nil {
			errorMsg := fmt.Sprintf("[RPC] In restarting, start container failed, containerID: %s", pod.Spec.Containers[i].ContainerID)
			log.ErrorLog(errorMsg)
			return err
		}
	}
	for i := 0; i < len(pod.Spec.Containers); i += 1 {
		pod.Spec.Containers[i].ContainerStatus = apiObject.Container_Running
	}

	return nil
}

func (r *RuntimeManager) StopPod(pod *apiObject.Pod) error {
	log.InfoLog("[RPC] Start StopPod")
	for i := 0; i < len(pod.Spec.Containers); i += 1 {
		_, err := r.runtimeClient.StopContainer(context.Background(), &runtimeapi.StopContainerRequest{
			ContainerId: pod.Spec.Containers[i].ContainerID,
			Timeout:     config.RPCRequestTimeout,
		})
		if err != nil {
			errorMsg := fmt.Sprintf("[RPC] Stop container failed, containerID: %s", pod.Spec.Containers[i].ContainerID)
			log.ErrorLog(errorMsg)
			return err
		}

		pod.Spec.Containers[i].ContainerStatus = apiObject.Container_Paused
	}
	return nil
}

func (r *RuntimeManager) DeletePod(pod *apiObject.Pod) error {
	log.InfoLog("[RPC] Start DeletePod")
	for i := 0; i < len(pod.Spec.Containers); i += 1 {
		_, err := r.runtimeClient.RemoveContainer(context.Background(), &runtimeapi.RemoveContainerRequest{
			ContainerId: pod.Spec.Containers[i].ContainerID,
		})
		if err != nil {
			errorMsg := fmt.Sprintf("[RPC] Remove container failed, containerID: %s", pod.Spec.Containers[i].ContainerID)
			log.ErrorLog(errorMsg)
			return err
		}

		pod.Spec.Containers[i].ContainerStatus = apiObject.Container_Removing
	}
	return nil
}

// 此处保留所有podSandbox，创建pod内部所有的容器
func (r *RuntimeManager) RecreatePodContainers(pod *apiObject.Pod) error {
	log.InfoLog("[RPC] Start RecreatePodContainers")

	sandboxConfig, err := r.getPodSandBoxConfig(pod)
	if err != nil {
		return err
	}
	containers := &pod.Spec.Containers
	for i := 0; i < len(*containers); i += 1 {
		containerConfig, err := r.getContainerConfig(&(*containers)[i], sandboxConfig)
		if err != nil {
			log.ErrorLog("generate container config failed")
			return err
		}

		containerID, err := r.CreateContainers(pod.PodSanboxId, containerConfig, sandboxConfig)
		if err != nil {
			log.ErrorLog("Create containers failed")
			return err
		}

		(*containers)[i].ContainerID = containerID
		(*containers)[i].ContainerStatus = apiObject.Container_Created
		message := fmt.Sprintf("container Id:%s is created ", containerID)
		log.InfoLog(message)
	}
	return nil
}

// TODO:似乎是在某个容器内部执行某条指令，不知道在哪里会被用到
func (r *RuntimeManager) ExecPodContainer(pod *apiObject.Pod) error {

	return nil
}
