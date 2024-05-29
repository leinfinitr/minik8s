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
	// 检查pv是否已经存在
	key := config.EtcdPvPrefix + "/" + pvNamespace + "/" + pvName
	response, _ := etcdclient.EtcdStore.Get(key)
	if response != "" {
		log.ErrorLog("Create PersistentVolume: pv already exists" + response)
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "pv already exists"})
		return
	}
	log.DebugLog("CreatePv: " + key)
	// 创建pv
	err = specctlrs.PvControllerInstance.AddPv(&pv)
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 将pv存入etcd
	pvJson, err := json.Marshal(&pv)
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(key, string(pvJson))
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": "Create PersistentVolume " + pvName})
}
