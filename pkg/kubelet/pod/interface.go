// 描述：此处的函数是供 kubelet 调用的接口

package pod

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"

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
	param := c.Param("param")
	exec := new(apiObject.ExecReq)

	exec.ContainerId = containerId
	exec.Cmd = []string{param}
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
