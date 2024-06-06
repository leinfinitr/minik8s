package conversion

import (
	"strconv"

	"minik8s/pkg/apiObject"
	"minik8s/tools/log"

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
			Containers: []apiObject.Container{
				{
					Name:  serverless.Name,
					Image: serverless.Image,
					Command: []string{
						"/bin/sh",
						"-c",
						"while true; do sleep 1000; done",
					},
					WorkingDir: "/mnt",
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

// ResourcesConvert 将一个 Resources 转化为单位为 KB 的大小
func ResourcesConvert(resources string) int {
	// resources 中第一个非数字字符之后的部分都是单位
	unit := ""
	number := ""
	for _, r := range resources {
		if r >= '0' && r <= '9' {
			number += string(r)
		} else {
			unit += string(r)
		}
	}
	// 将资源转化为 KB
	switch unit {
	case "Ki":
		numberInt, _ := strconv.Atoi(number)
		return 1 * numberInt
	case "Mi":
		numberInt, _ := strconv.Atoi(number)
		return 1024 * numberInt
	case "Gi":
		numberInt, _ := strconv.Atoi(number)
		return 1024 * 1024 * numberInt
	case "Ti":
		numberInt, _ := strconv.Atoi(number)
		return 1024 * 1024 * 1024 * numberInt
	default:
		log.ErrorLog("Unknown unit: " + unit)
		return 0
	}
}
