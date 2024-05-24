// 描述：此处的函数是供 kubelet 调用的接口

package pod

import (
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"
	"time"

	"github.com/gin-gonic/gin"
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
	}
}

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

// func StartPod(c *gin.Context) {
// 	var pod apiObject.Pod
// 	err := c.ShouldBindJSON(&pod)
// 	if err != nil {
// 		log.ErrorLog("StartPod error: " + err.Error())
// 	}
// 	err = podManager.StartPod(&pod)
// 	if err != nil {
// 		log.ErrorLog("StartPod error: " + err.Error())
// 	} else {
// 		c.JSON(200, "")
// 	}
// }

// func StopPod(c *gin.Context) {
// 	var pod apiObject.Pod
// 	err := c.ShouldBindJSON(&pod)
// 	if err != nil {
// 		log.ErrorLog("StopPod error: " + err.Error())
// 	}
// 	err = podManager.StopPod(&pod)
// 	if err != nil {
// 		log.ErrorLog("StopPod error: " + err.Error())
// 	} else {
// 		c.JSON(200, "")
// 	}
// }

// func RestartPod(c *gin.Context) {
// 	var pod apiObject.Pod
// 	err := c.ShouldBindJSON(&pod)
// 	if err != nil {
// 		log.ErrorLog("RestartPod error: " + err.Error())
// 	}
// 	err = podManager.RestartPod(&pod)
// 	if err != nil {
// 		log.ErrorLog("RestartPod error: " + err.Error())
// 	} else {
// 		c.JSON(200, "")
// 	}
// }

// func DeletePodByUUID(c *gin.Context) {
// 	var pod apiObject.Pod
// 	err := c.ShouldBindJSON(&pod)
// 	if err != nil {
// 		log.ErrorLog("DeletePodByUUID error: " + err.Error())
// 	}
// 	err = podManager.DeletePodByUUID(&pod)
// 	if err != nil {
// 		log.ErrorLog("DeletePodByUUID error: " + err.Error())
// 	} else {
// 		c.JSON(200, "")
// 	}
// }

// func RecreatePodContainer(c *gin.Context) {
// 	var pod apiObject.Pod
// 	err := c.ShouldBindJSON(&pod)
// 	if err != nil {
// 		log.ErrorLog("RecreatePodContainer error: " + err.Error())
// 	}
// 	err = podManager.RecreatePodContainer(&pod)
// 	if err != nil {
// 		log.ErrorLog("RecreatePodContainer error: " + err.Error())
// 	} else {
// 		c.JSON(200, "")
// 	}
// }

// func ExecPodContainer(c *gin.Context) {
// 	var pod apiObject.Pod
// 	err := c.ShouldBindJSON(&pod)
// 	if err != nil {
// 		log.ErrorLog("ExecPodContainer error: " + err.Error())
// 	}
// 	err = podManager.ExecPodContainer(&pod)
// 	if err != nil {
// 		log.ErrorLog("ExecPodContainer error: " + err.Error())
// 	} else {
// 		c.JSON(200, "")
// 	}
// }

// ScanPodStatus 用于扫描 pod 的状态，根据 pod 的状态进行相应的操作
func ScanPodStatus() {
	log.DebugLog("start scan pod status")
	// 每间隔 10s 扫描一次
	for {
		// 更新 pod 的状态
		err := podManager.UpdatePodStatus()
		if err != nil {
			log.ErrorLog("ScanPodStatus error: " + err.Error())
		}
		// 遍历所有的pod
		for _, pod := range podManager.PodMapByUUID {
			// 根据每个pod当前所处的阶段进行相应的操作
			phase := pod.Status.Phase
			switch phase {
			case apiObject.PodSucceeded:
				// 如果 pod 处于 Succeeded 阶段，则运行 pod
				go func() {
					err := podManager.StartPod(pod)
					if err != nil {
						log.ErrorLog("StartPod error: " + err.Error())
					}
				}()
				// 其余情况暂不处理
			default:
				log.DebugLog("Pod is in phase: " + phase)

			}
		}
		// 间隔10s
		time.Sleep(10 * time.Second)
	}
}

// ExecPodContainer 用于执行 pod 中的容器
func ExecPodContainer(c *gin.Context) {
	containerId := c.Param("container")
	param := c.Param("param")
	exec := new(apiObject.ExecReq)

	exec.ContainerId = containerId
	exec.Cmd = []string{param}
	exec.Tty = true
	exec.Stdin = true
	exec.Stdout = true
	exec.Stderr = true

	rep, err := podManager.ExecPodContainer(exec)
	if err != nil {
		log.ErrorLog("ExecPodContainer error: " + err.Error())
		c.JSON(config.HttpErrorCode, err.Error())
	} else {
		c.JSON(200, rep.Url)
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
