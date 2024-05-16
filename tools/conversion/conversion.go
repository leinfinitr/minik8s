package conversion

import (
	"minik8s/pkg/apiObject"
)

// ServerlessToPod 将一组 ServerlessFunction 转换为一组 Container
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
