package handlers

import (
	"github.com/gin-gonic/gin"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/controller"
	"minik8s/tools/log"
)

// CreatePVC 创建PersistentVolumeClaim
func CreatePVC(c *gin.Context) {
	pvc := &apiObject.PersistentVolumeClaim{}
	err := c.ShouldBindJSON(pvc)
	if err != nil {
		log.ErrorLog("CreatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	pvcName := pvc.Metadata.Name
	pvcNamespace := pvc.Metadata.Namespace
	if pvcName == "" || pvcNamespace == "" {
		log.ErrorLog("Create PersistentVolumeClaim: name or namespace is empty")
		c.JSON(400, gin.H{"error": "name or namespace is empty"})
		return
	}

	// 创建pvc
	log.DebugLog("CreatePvc: " + pvcNamespace + "/" + pvcName)
	err = controller.ControllerManagerInstance.AddPvc(pvc)
	if err != nil {
		log.ErrorLog("Create PersistentVolumeClaim: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": "Create PersistentVolumeClaim " + pvcName})
}
