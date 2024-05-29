package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"

	etcdclient "minik8s/pkg/apiServer/etcdClient"
	specctlrs "minik8s/pkg/controller/specCtlrs"
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

	key := config.EtcdPvcPrefix + "/" + pvcNamespace + "/" + pvcName
	response, _ := etcdclient.EtcdStore.Get(key)
	if response != "" {
		log.ErrorLog("Create PersistentVolumeClaim: pvc already exists" + response)
		c.JSON(400, gin.H{"error": "pvc already exists"})
		return
	}
	log.DebugLog("CreatePvc: " + pvcName)

	err = specctlrs.PvControllerInstance.AddPvc(pvc)
	if err != nil {
		log.ErrorLog("Create PersistentVolumeClaim: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pvcJson, err := json.Marshal(pvc)
	if err != nil {
		log.ErrorLog("Create PersistentVolumeClaim: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}
	err = etcdclient.EtcdStore.Put(key, string(pvcJson))
	if err != nil {
		log.ErrorLog("Create PersistentVolumeClaim: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": "Create PersistentVolumeClaim " + pvcName})
}
