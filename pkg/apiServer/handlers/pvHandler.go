package handlers

import (
	"github.com/gin-gonic/gin"
	"minik8s/pkg/apiObject"
	etcdclient "minik8s/pkg/apiServer/etcdClient"
	"minik8s/pkg/config"
	"minik8s/tools/log"
)

// CreatePV 创建PersistentVolume
func CreatePV(c *gin.Context) {
	pv := &apiObject.PersistentVolume{}
	err := c.ShouldBindJSON(pv)
	if err != nil {
		log.ErrorLog("CreatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	newPvName := pv.Metadata.Name
	if newPvName == "" {
		log.ErrorLog("Create PersistentVolume: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	key := config.EtcdPvPrefix + "/" + newPvName
	response, _ := etcdclient.EtcdStore.Get(key)
	if response != "" {
		log.ErrorLog("Create PersistentVolume: pv already exists" + response)
		c.JSON(400, gin.H{"error": "pv already exists"})
		return
	}
	log.InfoLog("CreatePv: " + newPvName)

	// 优先创建本地pv
	if pv.Spec.Local.Path != "" {
		// TODO: 创建本地挂载目录
	} else {
		// TODO: 创建远程挂载目录
	}

	c.JSON(200, gin.H{"data": "Create PersistentVolume " + newPvName})
}
