// 描述：此处的函数是供 kubelet 调用的接口

package pod

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"
	"minik8s/tools/mount"

	etcdclient "minik8s/pkg/apiServer/etcdClient"
	httprequest "minik8s/tools/httpRequest"
)

// UpdatePod 删除原有 pod 之后根据 pod 信息创建新的 pod
func UpdatePod(c *gin.Context) {
	var pod apiObject.Pod
	err := c.ShouldBindJSON(&pod)
	if err != nil {
		log.ErrorLog("UpdatePod error: " + err.Error())
	}

	err = podManager.DeletePod(&pod)
	if err != nil {
		log.ErrorLog("UpdatePod error: " + err.Error())
		c.JSON(config.HttpErrorCode, err.Error())
	}

	err = podManager.AddPod(&pod)
	if err != nil {
		log.ErrorLog("UpdatePod error: " + err.Error())
		c.JSON(config.HttpErrorCode, err.Error())
	} else {
		c.JSON(200, "")
	}
}

// CreatePod 创建 pod
func CreatePod(c *gin.Context) {
	var pod apiObject.Pod
	err := c.ShouldBindJSON(&pod)
	if err != nil {
		log.ErrorLog("AddPod error: " + err.Error())
		c.JSON(config.HttpErrorCode, err.Error())
		return
	}

	// 若pod使用了volume，则对其进行处理
	if pod.Spec.Volumes != nil {
		for _, volume := range pod.Spec.Volumes {
			// 处理使用了pvc的volume
			if volume.PersistentVolumeClaim != nil {
				pvcKey := config.EtcdPvcPrefix + "/" + pod.Metadata.Namespace + "/" + volume.PersistentVolumeClaim.ClaimName
				pvcResponse, _ := etcdclient.EtcdStore.Get(pvcKey)
				// 若pvc不存在，则返回错误
				if pvcResponse == "" {
					log.ErrorLog("CreatePod: PVC not found")
					c.JSON(400, gin.H{"error": "PVC not found"})
					return
				}
				// 若pvc存在，则检查pvc的状态是否为Bound、是否已经绑定到pod
				pvc := &apiObject.PersistentVolumeClaim{}
				err = json.Unmarshal([]byte(pvcResponse), pvc)
				if err != nil {
					log.ErrorLog("CreatePod: " + err.Error())
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				if pvc.Status.Phase != apiObject.ClaimBound || pvc.Status.IsBound {
					log.ErrorLog("CreatePod: PVC can't be used")
					c.JSON(400, gin.H{"error": "PVC can't be used"})
					return
				}
				// 将pod绑定到pvc
				url := config.PVServerURL() + config.PersistentVolumeClaimURI
				url = strings.Replace(url, config.NameSpaceReplace, pod.Metadata.Namespace, -1)
				url = strings.Replace(url, config.NameReplace, pod.Metadata.Name, -1)
				res, err := httprequest.PostObjMsg(url, pvc)
				if err != nil {
					log.ErrorLog("Could not post the object message." + err.Error())
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
				if res == nil {
					log.ErrorLog("CreatePod: res is nil")
					c.JSON(500, gin.H{"error": "res is nil"})
					return
				}
				if res.StatusCode != http.StatusOK {
					log.ErrorLog("CreatePod: " + res.Status)
					c.JSON(res.StatusCode, gin.H{"error": res.Status})
					return
				}
				// 获取pvKey
				url = config.PVServerURL() + config.PersistentVolumeClaimURI
				url = strings.Replace(url, config.NameSpaceReplace, pvc.Metadata.Namespace, -1)
				url = strings.Replace(url, config.NameReplace, pvc.Metadata.Name, -1)
				res, err = httprequest.GetMsg(url)
				if err != nil {
					log.ErrorLog("Could not post the object message." + err.Error())
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
				if res == nil {
					log.ErrorLog("CreatePod: res is nil")
					c.JSON(http.StatusInternalServerError, gin.H{"error": "res is nil"})
					return
				}
				if res.StatusCode != http.StatusOK {
					log.ErrorLog("CreatePod: " + res.Status)
					c.JSON(res.StatusCode, gin.H{"error": res.Status})
					return
				}
				// 解析pvKey
				var pvKey string
				err = json.NewDecoder(res.Body).Decode(&pvKey)
				if err != nil {
					log.ErrorLog("CreatePod: " + err.Error())
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				if pvKey == "" {
					log.ErrorLog("CreatePod: pvName is empty")
					c.JSON(400, gin.H{"error": "pvName is empty"})
					return
				}
				log.DebugLog("Parse pv name: " + pvKey)
				// 将本地挂载目录挂载到服务器
				err = mount.LocalToServer(pvKey)
				if err != nil {
					log.ErrorLog("CreatePod: " + err.Error())
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				// 为pod中所有使用该持久化卷挂载的容器添加Mount
				mount.AddMountsToContainer(&pod, volume, config.PVClientPath+"/"+pvKey)
			}
			// 处理使用了emptyDir的volume
			if volume.EmptyDir.SizeLimit != "" {
				// 创建emptyDir
				cleanCmd := "rm -rf " + config.DefaultVolumePath + "/" + pod.Metadata.Name + "/" + volume.Name
				mkdirCmd := "mkdir -p " + config.DefaultVolumePath + "/" + pod.Metadata.Name + "/" + volume.Name
				cmd := exec.Command("sh", "-c", cleanCmd)
				err = cmd.Run()
				if err != nil {
					log.ErrorLog("CreatePod: " + err.Error())
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				cmd = exec.Command("sh", "-c", mkdirCmd)
				err = cmd.Run()
				if err != nil {
					log.ErrorLog("CreatePod: " + err.Error())
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				mount.AddMountsToContainer(&pod, volume, config.DefaultVolumePath+"/"+pod.Metadata.Name+"/"+volume.Name)
			}
			// 处理使用了hostPath的volume
			if volume.HostPath.Path != "" {
				// 如果hostPath不存在，则返回错误
				if _, err := os.Stat(volume.HostPath.Path); os.IsNotExist(err) {
					log.ErrorLog("CreatePod: hostPath not found")
					c.JSON(400, gin.H{"error": "hostPath not found"})
					return
				}
				mount.AddMountsToContainer(&pod, volume, volume.HostPath.Path)
			}
		}
	}

	err = podManager.AddPod(&pod)
	if err != nil {
		log.ErrorLog("AddPod error: " + err.Error())
		c.JSON(config.HttpErrorCode, err.Error())
		return
	} else {
		c.JSON(config.HttpSuccessCode, pod)
	}
}

// DeletePod 删除 pod
func DeletePod(c *gin.Context) {
	var pod apiObject.Pod
	err := c.ShouldBindJSON(&pod)
	if err != nil {
		log.ErrorLog("DeletePod error: " + err.Error())
	}
	err = podManager.DeletePod(&pod)
	if err != nil {
		log.ErrorLog("DeletePod error: " + err.Error())
	} else {
		c.JSON(200, "")
	}
}

// GetPodStatus 获取 pod 的状态
func GetPodStatus(c *gin.Context) {
	var pod apiObject.Pod
	err := c.ShouldBindJSON(&pod)
	if err != nil {
		log.ErrorLog("GetPodStatus error: " + err.Error())
	}

	err = podManager.UpdatePodStatus()
	if err != nil {
		log.ErrorLog("GetPodStatus error: " + err.Error())
	}

	// 遍历 podManager 中的 PodMapByUUID，找到对应的 pod，返回 pod 的状态
	for _, v := range podManager.PodMapByUUID {
		if v.Metadata.UUID == pod.Metadata.UUID {
			c.JSON(200, v.Status)
			return
		}
	}

	// 如果没有找到对应的 pod，返回错误信息
	c.JSON(config.HttpErrorCode, "Pod not found")
}

// SyncPods 用于同步 pod
func SyncPods(c *gin.Context) {
	var pods []apiObject.Pod
	err := c.ShouldBindJSON(&pods)
	if err != nil {
		log.ErrorLog("SyncPod error: " + err.Error())
	}
	err = podManager.SyncPods(&pods)
	if err != nil {
		log.ErrorLog("SyncPod error: " + err.Error())
	} else {
		// 主动调用一次扫描 pod 状态的函数
		ScanPodStatusRoutine()
		c.JSON(200, "")
	}
}

// ScanPodStatus 用于扫描 pod 的状态，根据 pod 的状态进行相应的操作
func ScanPodStatus() {
	log.InfoLog("start scan pod status")
	// 每间隔 15s 扫描一次
	for {
		ScanPodStatusRoutine()
		time.Sleep(15 * time.Second)
	}
}

// ScanPodStatusRoutine 用于扫描 pod 的状态，根据 pod 的状态进行相应的操作
func ScanPodStatusRoutine() {
	// 更新 pod 的状态
	err := podManager.UpdatePodStatus()
	if err != nil {
		log.ErrorLog("ScanPodStatus error: " + err.Error())
	}
	// 遍历所有的pod，根据pod的状态进行相应的操作
	for _, pod := range podManager.PodMapByUUID {
		// 防止协程中因为和主协程共享变量的变化
		pod := pod
		// 根据每个pod当前所处的阶段进行相应的操作
		phase := pod.Status.Phase
		switch phase {
		case apiObject.PodSucceeded:
			// 如果 pod 处于 Succeeded 阶段，则运行 pod
			go func() {
				err := podManager.StartPod(pod)
				if err != nil {
					log.ErrorLog("StartPod error: " + err.Error())
				} else {
					// 更新apiServer的pod状态
					UpdatePodStatus(pod)
				}
			}()
		default:
			// 其余情况暂不处理
			log.DebugLog("Pod " + pod.Metadata.Name + " is in phase: " + string(phase))
		}
	}
}

// ExecPodContainer 用于在 pod 中的容器执行命令
func ExecPodContainer(c *gin.Context) {
	containerId := c.Param("container")
	var cmd apiObject.Command
	exec := new(apiObject.ExecReq)

	if err := c.ShouldBindJSON(&cmd); err != nil {
		log.ErrorLog("ExecPodContainer error: " + err.Error())
		c.JSON(config.HttpErrorCode, err.Error())
		return
	}
	exec.ContainerId = containerId
	exec.Cmd = []string{cmd.Cmd}
	exec.Tty = true
	exec.Stdin = true
	exec.Stdout = true
	exec.Stderr = false

	rep, err := podManager.ExecPodContainer(exec)
	if err != nil {
		log.ErrorLog("ExecPodContainer error: " + err.Error())
		c.JSON(config.HttpErrorCode, err.Error())
	} else {
		log.DebugLog("ExecPodContainer success: " + rep)
		c.JSON(config.HttpSuccessCode, rep)
	}
}

// GetPods 用于获取所有的 pod
func GetPods(c *gin.Context) {
	log.DebugLog("GetPods")
	pods := new([]apiObject.Pod)
	for _, pod := range podManager.PodMapByUUID {
		*pods = append(*pods, *pod)
	}
	c.JSON(200, pods)
}

// UpdatePodStatus 向 apiServer 更新 pod 的状态
func UpdatePodStatus(pod *apiObject.Pod) {
	// 不断重试，直到更新成功
	for {
		url := config.HttpSchema + config.APIServerLocalAddress + ":" + fmt.Sprint(config.APIServerLocalPort) + config.PodStatusURI
		url = strings.Replace(url, config.NameSpaceReplace, pod.Metadata.Namespace, -1)
		url = strings.Replace(url, config.NameReplace, pod.Metadata.Name, -1)
		res, err := httprequest.PutObjMsg(url, pod.Status)
		if err != nil {
			log.ErrorLog("UpdatePodStatus: " + err.Error())
		} else if res.StatusCode != 200 {
			log.ErrorLog("UpdatePodStatus error: " + res.Status)
		} else {
			break
		}
	}
}
