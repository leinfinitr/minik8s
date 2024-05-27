package conversion

import (
	"minik8s/pkg/apiObject"

	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

// ServerlessToPod 将一个 Serverless 对象转换为 Pod 对象
func ServerlessToPod(serverless apiObject.Serverless) apiObject.Pod {
	pod := apiObject.Pod{
		TypeMeta: apiObject.TypeMeta{
			Kind:       apiObject.PodType,
			APIVersion: "v1",
		},
		Metadata: apiObject.ObjectMeta{
			Name:      serverless.Name,
			Namespace: "serverless",
		},
		Spec: apiObject.PodSpec{
			Volumes: []apiObject.Volume{
				{
					Name: serverless.Volume,
				},
			},
			Containers: []apiObject.Container{
				{
					Name:       serverless.Name,
					Image:      serverless.Image,
					Command:    []string{"pip", "install", "-r", "/mnt/requirements.txt"},
					WorkingDir: "/mnt",
					VolumeMounts: []apiObject.VolumeMount{
						{
							Name:      serverless.Volume,
							MountPath: "/mnt",
						},
					},
					Mounts: []*apiObject.Mount{
						{
							HostPath:      serverless.HostPath,
							ContainerPath: "/mnt",
							ReadOnly:      false,
						},
					},
				},
			},
		},
	}
	return pod
}

// PodToServerless 将一个 Pod 对象转换为 Serverless 对象
func PodToServerless(pod apiObject.Pod) apiObject.Serverless {
	serverless := apiObject.Serverless{
		Name: pod.Metadata.Name,
	}
	for _, container := range pod.Spec.Containers {
		serverless.Image = container.Image
		serverless.Volume = container.VolumeMounts[0].Name
	}
	return serverless
}

// MountsToMounts 将一个 Mount 对象数组转换为 ContainerConfig.Mounts 对象数组
func MountsToMounts(mounts []*apiObject.Mount) []*runtimeapi.Mount {
	configMounts := make([]*runtimeapi.Mount, 0)
	for _, mount := range mounts {
		configMount := &runtimeapi.Mount{
			HostPath:      mount.HostPath,
			ContainerPath: mount.ContainerPath,
			Readonly:      mount.ReadOnly,
		}
		configMounts = append(configMounts, configMount)
	}
	return configMounts
}
