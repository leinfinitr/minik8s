package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"

	httprequest "minik8s/tools/httpRequest"
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

	// 转发给pvController
	log.DebugLog("CreatePvc: " + pvcNamespace + "/" + pvcName)
	url := config.PVServerURL() + config.PersistentVolumeClaimsURI
	res, err := httprequest.PostObjMsg(url, pvc)
	if err != nil {
		log.ErrorLog("Could not post the object message." + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if res.StatusCode != http.StatusCreated {
		log.ErrorLog("Create PersistentVolumeClaim: " + res.Status)
		c.JSON(500, gin.H{"error": res.Status})
		return
	}

	c.JSON(200, gin.H{"data": "Create PersistentVolumeClaim " + pvcName})
}
