package mount

import (
	"os/exec"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"
)

// AddMountsToContainer 将一个 Mount 对象添加到 Container.Mounts
func AddMountsToContainer(pod *apiObject.Pod, volume apiObject.Volume, hostPath string) {
	for i, container := range pod.Spec.Containers {
		for _, volumeMount := range container.VolumeMounts {
			if volumeMount.Name == volume.Name {
				// 为容器添加Mount
				if container.Mounts == nil {
					container.Mounts = make([]*apiObject.Mount, 0)
				}
				mount := &apiObject.Mount{
					HostPath:      hostPath,
					ContainerPath: volumeMount.MountPath,
					ReadOnly:      false,
				}
				container.Mounts = append(container.Mounts, mount)
			}
		}
		// 替换pod中的容器
		pod.Spec.Containers[i] = container
	}
}

// LocalToServer 将本地目录挂载到服务器
func LocalToServer(localPath string) error {
	// 将本地目录 /pvclient 挂载到服务器目录 /pvserver
	mountCmd := "mount " + config.NFSServer + ":" + config.PVServerPath + " " + config.PVClientPath
	cmd := exec.Command("sh", "-c", mountCmd)
	err := cmd.Run()
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		return err
	}
	log.DebugLog("Bind to NFS server: " + config.NFSServer + ":" + config.PVServerPath)
	// 在目录 /pvclient 创建目录 /:namespace/:name 作为PersistentVolume
	mkdirCmd := "mkdir -p " + localPath
	cmd = exec.Command("sh", "-c", mkdirCmd)
	err = cmd.Run()
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		return err
	}
	log.DebugLog("Create PersistentVolume: " + localPath)
	// 清空目录 /pvclient/:namespace/:name
	rmCmd := "rm -rf " + localPath + "/*"
	cmd = exec.Command("sh", "-c", rmCmd)
	err = cmd.Run()
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		return err
	}
	log.DebugLog("Clean PersistentVolume: " + localPath)
	return nil
}
