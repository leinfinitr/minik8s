// 描述：此处的函数是供 kubelet 中 KubeletAPIRouter 调用的接口
// 每个函数的功能包括：
// 1. 对 http 请求中的参数进行处理解析
// 2. 转发给 podManager

package pod

import (
	"github.com/gin-gonic/gin"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"
)

func AddPod(c *gin.Context) {
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
	} else {
		c.JSON(200, "")
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

func StartPod(c *gin.Context) {
	var pod apiObject.Pod
	err := c.ShouldBindJSON(&pod)
	if err != nil {
		log.ErrorLog("StartPod error: " + err.Error())
	}
	err = podManager.StartPod(&pod)
	if err != nil {
		log.ErrorLog("StartPod error: " + err.Error())
	} else {
		c.JSON(200, "")
	}
}

func StopPod(c *gin.Context) {
	var pod apiObject.Pod
	err := c.ShouldBindJSON(&pod)
	if err != nil {
		log.ErrorLog("StopPod error: " + err.Error())
	}
	err = podManager.StopPod(&pod)
	if err != nil {
		log.ErrorLog("StopPod error: " + err.Error())
	} else {
		c.JSON(200, "")
	}
}

func RestartPod(c *gin.Context) {
	var pod apiObject.Pod
	err := c.ShouldBindJSON(&pod)
	if err != nil {
		log.ErrorLog("RestartPod error: " + err.Error())
	}
	err = podManager.RestartPod(&pod)
	if err != nil {
		log.ErrorLog("RestartPod error: " + err.Error())
	} else {
		c.JSON(200, "")
	}
}

func DeletePodByUUID(c *gin.Context) {
	var pod apiObject.Pod
	err := c.ShouldBindJSON(&pod)
	if err != nil {
		log.ErrorLog("DeletePodByUUID error: " + err.Error())
	}
	err = podManager.DeletePodByUUID(&pod)
	if err != nil {
		log.ErrorLog("DeletePodByUUID error: " + err.Error())
	} else {
		c.JSON(200, "")
	}
}

func RecreatePodContainer(c *gin.Context) {
	var pod apiObject.Pod
	err := c.ShouldBindJSON(&pod)
	if err != nil {
		log.ErrorLog("RecreatePodContainer error: " + err.Error())
	}
	err = podManager.RecreatePodContainer(&pod)
	if err != nil {
		log.ErrorLog("RecreatePodContainer error: " + err.Error())
	} else {
		c.JSON(200, "")
	}
}

func ExecPodContainer(c *gin.Context) {
	var pod apiObject.Pod
	err := c.ShouldBindJSON(&pod)
	if err != nil {
		log.ErrorLog("ExecPodContainer error: " + err.Error())
	}
	err = podManager.ExecPodContainer(&pod)
	if err != nil {
		log.ErrorLog("ExecPodContainer error: " + err.Error())
	} else {
		c.JSON(200, "")
	}
}
