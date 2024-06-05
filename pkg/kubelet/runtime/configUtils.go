package runtime

import (
	"minik8s/pkg/apiObject"
	"minik8s/tools/conversion"
	"minik8s/tools/log"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func (r *RuntimeManager) getPodSandBoxConfig(pod *apiObject.Pod) (*runtimeapi.PodSandboxConfig, error) {
	// put basic infos from pod into config
	logDirectory := "/var/log/pods" + pod.Metadata.Name

	podSandboxConfig := &runtimeapi.PodSandboxConfig{
		Metadata: &runtimeapi.PodSandboxMetadata{
			Name:      pod.Metadata.Name,
			Namespace: pod.Metadata.Namespace,
			Uid:       pod.Metadata.UUID,
			Attempt:   1,
		},
		Labels:       pod.Metadata.Labels,
		Annotations:  pod.Metadata.Annotations,
		LogDirectory: logDirectory, //`/var/log/pods/<NAMESPACE>_<NAME>_<UID>/`
	}

	// TODO: 需要获取config中基本的DNS信息，暂时不需要

	podSandboxConfig.Hostname = pod.Spec.NodeName
	podSandboxConfig.LogDirectory = "/var/log/pods" //`/var/log/pods/<NAMESPACE>_<NAME>_<UID>/`

	// TODO: 这里可能还需要实现端口映射

	// TODO: 默认需要生成关于linux的配置
	linuxConfig, err := r.getPodSandBoxLinuxConfig()
	if err != nil {
		return nil, nil
	}
	podSandboxConfig.Linux = linuxConfig

	return podSandboxConfig, nil
}

func (r *RuntimeManager) getPodSandBoxLinuxConfig() (*runtimeapi.LinuxPodSandboxConfig, error) {
	linuxConfig := &runtimeapi.LinuxPodSandboxConfig{
		CgroupParent: "",
		SecurityContext: &runtimeapi.LinuxSandboxSecurityContext{
			Privileged: false,
			Seccomp: &runtimeapi.SecurityProfile{
				ProfileType: runtimeapi.SecurityProfile_RuntimeDefault,
			},
		},
	}

	// 我们默认pod关于安全上下文的配置均为空，跳过该部分的配置

	return linuxConfig, nil
}

// 生成 ContainerConfig，可供runtimeClient直接使用发送
func (r *RuntimeManager) getContainerConfig(container *apiObject.Container, sandboxConfig *runtimeapi.PodSandboxConfig) (*runtimeapi.ContainerConfig, error) {
	// 1. 将镜像拉取到本地
	imageRef, err := r.imageManager.PullImage(container, sandboxConfig)
	if err != nil {
		return nil, err
	}

	logPath := container.Name + ".log"
	// 2. 创建container
	// 需要在dns中追加/etc/hosts文件，以便容器内部可以访问到集群内部署的DNS服务
	container.Mounts = append(container.Mounts, &apiObject.Mount{
		ContainerPath: "/etc/hosts",
		HostPath:      "/etc/hosts",
		ReadOnly:      true,
	})

	config := &runtimeapi.ContainerConfig{
		Metadata: &runtimeapi.ContainerMetadata{
			Name:    container.Name,
			Attempt: 1,
		},
		Image: &runtimeapi.ImageSpec{
			Image:              imageRef,
			UserSpecifiedImage: container.Image,
		},
		Command:    container.Command,
		WorkingDir: container.WorkingDir,
		Mounts:     conversion.MountsToMounts(container.Mounts),
		LogPath:    logPath,
	}
	log.DebugLog("ContainerConfig: " + config.String())

	return config, nil

}
