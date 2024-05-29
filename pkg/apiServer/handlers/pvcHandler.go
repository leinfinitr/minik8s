package handlers

import (
	"github.com/gin-gonic/gin"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"

	etcdclient "minik8s/pkg/apiServer/etcdClient"
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
	newPvcName := pvc.Metadata.Name
	if newPvcName == "" {
		log.ErrorLog("Create PersistentVolumeClaim: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}

	key := config.EtcdPvcPrefix + "/" + newPvcName
	response, _ := etcdclient.EtcdStore.Get(key)
	if response != "" {
		log.ErrorLog("Create PersistentVolumeClaim: already exists" + response)
		c.JSON(400, gin.H{"error": "pvc already exists"})
		return
	}
	log.InfoLog("CreatePv: " + newPvcName)

	// TODO: 与PV绑定

	c.JSON(200, gin.H{"data": "Create PersistentVolumeClaim " + newPvcName})
}
