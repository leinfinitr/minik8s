package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"

	httprequest "minik8s/tools/httpRequest"
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
	url := config.PVServerURL() + config.PersistentVolumesURI
	res, err := httprequest.PostObjMsg(url, pv)
	if err != nil {
		log.ErrorLog("Could not post the object message." + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if res.StatusCode != http.StatusCreated {
		log.ErrorLog("Create PersistentVolume: " + res.Status)
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Status})
		return
	}

	c.JSON(200, gin.H{"data": "Create PersistentVolume " + pvName})
}
