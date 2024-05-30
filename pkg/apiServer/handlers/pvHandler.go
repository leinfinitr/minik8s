package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"minik8s/pkg/apiObject"
	"minik8s/tools/log"

	specctlrs "minik8s/pkg/controller/specCtlrs"
)

// CreatePV 创建PersistentVolume
func CreatePV(c *gin.Context) {
	var pv apiObject.PersistentVolume
	err := c.ShouldBindJSON(pv)
	if err != nil {
		log.ErrorLog("CreatePod: " + err.Error())
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	pvName := pv.Metadata.Name
	pvNamespace := pv.Metadata.Namespace
	if pvName == "" || pvNamespace == "" {
		log.ErrorLog("Create PersistentVolume: name or namespace is empty")
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "name or namespace is empty"})
		return
	}

	// 创建pv
	log.DebugLog("CreatePv: " + pvNamespace + "/" + pvName)
	err = specctlrs.PvControllerInstance.AddPv(&pv)
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": "Create PersistentVolume " + pvName})
}
