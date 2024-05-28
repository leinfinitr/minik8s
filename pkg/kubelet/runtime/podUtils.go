package runtime

import (
	"context"
	"errors"
	"fmt"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/host"
	"minik8s/tools/log"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

type PodUtils interface {
	CreatePod(pod *apiObject.Pod) error
	StartPod(pod *apiObject.Pod) error
	RestartPod(pod *apiObject.Pod) error
	StopPod(pod *apiObject.Pod) error
	DeletePod(pod *apiObject.Pod) error
	GetPodSandboxStatus(podId string) (*runtimeapi.PodSandboxStatus, error)
	RecreatePodContainer(pod *apiObject.Pod) error
	ExecPodContainer(req *apiObject.ExecReq) (*apiObject.ExecRsp, error)
	UpdatePodStatus(pod *apiObject.Pod) error
}

// CreatePod 在这里，我们创建一个Pod相当于是创建一个Sandbox，并且会创建Pod内部的所有容器
func (r *RuntimeManager) CreatePod(pod *apiObject.Pod) error {
	log.InfoLog("[RPC] Start CreatePod")

	pod.Status.Phase = apiObject.PodBuilding
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
	response, err := r.runtimeClient.RunPodSandbox(context.Background(), request)
	if err != nil {
		log.ErrorLog("Create Pod sandbox fail " + err.Error())
		return nil
	}

	log.DebugLog("Create Pod sandbox Success")

	if response.PodSandboxId == "" {
		errorMessage := "PodSandboxId set for pod sandbox is failed "
		log.ErrorLog(errorMessage)
		return errors.New(errorMessage)
	}
	message := fmt.Sprintf("PodSandboxId is set successfullly, Is : " + response.PodSandboxId)
	log.InfoLog(message)
	pod.PodSandboxId = response.PodSandboxId

	// 调用接口去获取分配给Pod的IP信息
	res, err := r.GetPodSandboxStatus(pod.PodSandboxId)
	if err != nil {
		log.ErrorLog("GetPodSandboxStatus failed")
		return err
	}
	pod.Status.PodIP = res.Network.Ip // 获取到Pod的IP信息

	// 调用接口去创建Pod内部的所有容器
	containers := &pod.Spec.Containers
	for i := 0; i < len(*containers); i += 1 {
		containerConfig, err := r.getContainerConfig(&(*containers)[i], sandboxConfig)
		if err != nil {
			log.ErrorLog("generate container config failed")
			return err
		}

		containerID, err := r.CreateContainers(pod.PodSandboxId, containerConfig, sandboxConfig)
		if err != nil {
			log.ErrorLog("Create containers failed")
			return err
		}

		(*containers)[i].ContainerID = containerID
		(*containers)[i].ContainerStatus = apiObject.ContainerCreated
		message := fmt.Sprintf("container Id:%s is created ", containerID)
		log.InfoLog(message)
	}

	pod.Status.Phase = apiObject.PodSucceeded

	return nil
}

// CreateContainers 创建指定配置文件的container
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

// StartPod 运行该Pod内部的所有容器
func (r *RuntimeManager) StartPod(pod *apiObject.Pod) error {
	log.InfoLog("[RPC] Start StartPod")
	// 运行所有的容器
	for i := 0; i < len(pod.Spec.Containers); i += 1 {
		_, err := r.runtimeClient.StartContainer(context.Background(), &runtimeapi.StartContainerRequest{
			ContainerId: pod.Spec.Containers[i].ContainerID,
		})
		if err != nil {
			errorMsg := fmt.Sprintf("[RPC] Start container failed, containerID: %s", pod.Spec.Containers[i].ContainerID)
			log.ErrorLog(errorMsg)
			return err
		}
		pod.Spec.Containers[i].ContainerStatus = apiObject.ContainerRunning

	}
	pod.Status.Phase = apiObject.PodRunning
	return nil
}

func (r *RuntimeManager) RestartPod(pod *apiObject.Pod) error {
	log.InfoLog("[RPC] Start RestartPod")
	// 考虑到容器之间可能存在依赖，为了保证可用性，在暂停所有的容器后再重新启动
	pod.Status.Phase = apiObject.PodBuilding
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
		pod.Spec.Containers[i].ContainerStatus = apiObject.ContainerRunning
	}

	pod.Status.Phase = apiObject.PodRunning
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

	}
	return nil
}

func (r *RuntimeManager) GetPodSandboxStatus(podId string) (*runtimeapi.PodSandboxStatus, error) {
	log.InfoLog("[RPC] Start GetPodSandboxStatus")
	response, err := r.runtimeClient.PodSandboxStatus(context.Background(), &runtimeapi.PodSandboxStatusRequest{
		PodSandboxId: podId,
	})

	if err != nil {
		log.ErrorLog("GetPodSandboxStatus failed")
		return nil, err
	}

	return response.Status, nil
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
	}
	return nil
}

// RecreatePodContainers 此处保留所有podSandbox，创建pod内部所有的容器
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

		containerID, err := r.CreateContainers(pod.PodSandboxId, containerConfig, sandboxConfig)
		if err != nil {
			log.ErrorLog("Create containers failed")
			return err
		}

		(*containers)[i].ContainerID = containerID
		(*containers)[i].ContainerStatus = apiObject.ContainerCreated
		message := fmt.Sprintf("container Id:%s is created ", containerID)
		log.InfoLog(message)
	}
	return nil
}

func (r *RuntimeManager) ExecPodContainer(req *apiObject.ExecReq) (*apiObject.ExecRsp, error) {
	log.InfoLog("[RPC] Start ExecPodContainer")

	response, err := r.runtimeClient.Exec(context.Background(), &runtimeapi.ExecRequest{
		ContainerId: req.ContainerId,
		Cmd:         req.Cmd,
		Tty:         req.Tty,
		Stdin:       req.Stdin,
		Stdout:      req.Stdout,
		Stderr:      req.Stderr,
	})

	if err != nil {
		log.ErrorLog("Exec container failed")
		return nil, err
	}

	return &apiObject.ExecRsp{
		Url: response.Url,
	}, nil
}

func (r *RuntimeManager) UpdatePodStatus(pod *apiObject.Pod) error {

	// 记录所有容器的资源占用情况
	var cpuUsage, memoryUsage float64

	memoryAll, err := host.GetTotalMemory()
	if err != nil {
		log.ErrorLog("GetTotalMemory failed" + err.Error())
		return err
	}

	for _, container := range pod.Spec.Containers {
		response1, err := r.runtimeClient.ContainerStats(context.Background(), &runtimeapi.ContainerStatsRequest{
			ContainerId: container.ContainerID,
		})

		if response1 != nil {
			log.InfoLog("CPU usage" + fmt.Sprint(response1.Stats.Cpu.UsageCoreNanoSeconds.Value))
			log.InfoLog("Memory usage" + fmt.Sprint(response1.Stats.Memory.UsageBytes.Value))
		}

		// if err != nil {
		// 	log.ErrorLog("Container status from CRI failed" + err.Error())
		// 	container.ContainerStatus = apiObject.ContainerUnknown
		// 	return err
		// }

		response2, err := r.runtimeClient.ContainerStatus(context.Background(), &runtimeapi.ContainerStatusRequest{
			ContainerId: container.ContainerID,
		})
		if err != nil {
			log.ErrorLog("Container status from CRI failed" + err.Error())
			container.ContainerStatus = apiObject.ContainerUnknown
			return err
		}

		if (uint64(response1.Stats.Cpu.Timestamp) - uint64(response2.Status.StartedAt)) != 0 {
			cpuUsage += float64(response1.Stats.Cpu.UsageCoreNanoSeconds.Value / (uint64(response1.Stats.Cpu.Timestamp) - uint64(response2.Status.StartedAt)))
		} else {
			cpuUsage += 0
		}

		if memoryAll != 0 {
			memoryUsage += float64(response1.Stats.Memory.UsageBytes.Value / memoryAll)
		} else {
			memoryUsage += 0
		}

		switch response2.Status.State {
		case runtimeapi.ContainerState_CONTAINER_CREATED:
			container.ContainerStatus = apiObject.ContainerCreated
			pod.Status.Phase = apiObject.PodSucceeded
		case runtimeapi.ContainerState_CONTAINER_RUNNING:
			container.ContainerStatus = apiObject.ContainerRunning
		default:
			container.ContainerStatus = apiObject.ContainerUnknown
			pod.Status.Phase = apiObject.PodFailed
			return errors.New("container status isn't normal")
		}

		pod.Status.CpuUsage = cpuUsage
		pod.Status.MemUsage = memoryUsage
	}

	return nil
}
