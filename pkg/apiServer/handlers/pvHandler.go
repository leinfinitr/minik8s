package handlers

import (
	"encoding/json"
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

// GetPV 获取PVC绑定的PVKey
func GetPV(c *gin.Context) {
	pvName := c.Param("name")
	pvNamespace := c.Param("namespace")
	if pvName == "" || pvNamespace == "" {
		log.ErrorLog("Get PersistentVolume: name or namespace is empty")
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "name or namespace is empty"})
		return
	}

	// 获取pv
	log.DebugLog("GetPv: " + pvNamespace + "/" + pvName)
	url := config.PVServerURL() + config.PersistentVolumesURI + "/" + pvNamespace + "/" + pvName
	res, err := httprequest.GetMsg(url)
	if err != nil {
		log.ErrorLog("Could not get the object message." + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if res.StatusCode != http.StatusOK {
		log.ErrorLog("Get PersistentVolume: " + res.Status)
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Status})
		return
	}

	// 解析pvKey
	var pvKey string
	err = json.NewDecoder(res.Body).Decode(&pvKey)
	if err != nil {
		log.ErrorLog("Get PersistentVolume: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if pvKey == "" {
		log.ErrorLog("Get PersistentVolume: pvName is empty")
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "pvName is empty"})
		return
	}
	log.DebugLog("Parse pv name: " + pvKey)

	c.JSON(200, pvKey)
}
