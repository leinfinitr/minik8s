package conversion

import (
	"minik8s/pkg/apiObject"
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
					Name:            serverless.Name,
					Image:           serverless.Image,
					ImagePullPolicy: "IfNotPresent",
					VolumeMounts: []apiObject.VolumeMount{
						{
							Name:      serverless.Volume,
							MountPath: "/mnt",
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
