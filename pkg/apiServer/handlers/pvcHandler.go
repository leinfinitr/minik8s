package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

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

// BindPVC 将Pod绑定到持久化卷声明
func BindPVC(c *gin.Context) {
	pvc := &apiObject.PersistentVolumeClaim{}
	err := c.ShouldBindJSON(pvc)
	if err != nil {
		log.ErrorLog("BindPVC: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	podName := c.Param("name")
	podNamespace := c.Param("namespace")
	if podName == "" || podNamespace == "" {
		log.ErrorLog("BindPVC: name or namespace is empty")
		c.JSON(400, gin.H{"error": "name or namespace is empty"})
		return
	}

	// 转发给pvController
	url := config.PVServerURL() + config.PersistentVolumeClaimURI
	url = strings.Replace(url, config.NameSpaceReplace, podNamespace, -1)
	url = strings.Replace(url, config.NameReplace, podName, -1)
	res, err := httprequest.PostObjMsg(url, pvc)
	if err != nil {
		log.ErrorLog("Could not post the object message." + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if res.StatusCode != http.StatusOK {
		log.ErrorLog("BindPVC: " + res.Status)
		c.JSON(500, gin.H{"error": res.Status})
		return
	}

	c.JSON(200, gin.H{"data": "Bind PersistentVolumeClaim " + pvc.Metadata.Name + " to Pod " + podName})
}

// GetPVC 获取指定持久化卷声明
func GetPVC(c *gin.Context) {
	pvcName := c.Param("name")
	pvcNamespace := c.Param("namespace")
	if pvcName == "" || pvcNamespace == "" {
		log.ErrorLog("GetPVC: name or namespace is empty")
		c.JSON(400, gin.H{"error": "name or namespace is empty"})
		return
	}

	// 转发给pvController
	url := config.PVServerURL() + config.PersistentVolumeClaimURI
	url = strings.Replace(url, config.NameSpaceReplace, pvcNamespace, -1)
	url = strings.Replace(url, config.NameReplace, pvcName, -1)
	res, err := httprequest.GetMsg(url)
	if err != nil {
		log.ErrorLog("Could not get the object message." + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if res.StatusCode != http.StatusOK {
		log.ErrorLog("GetPVC: " + res.Status)
		c.JSON(500, gin.H{"error": res.Status})
		return
	}

	var pvc apiObject.PersistentVolumeClaim
	err = json.NewDecoder(res.Body).Decode(&pvc)
	if err != nil {
		log.ErrorLog("GetPVC: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, pvc)
}
