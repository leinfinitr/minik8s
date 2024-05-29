package handlers

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"

	etcdclient "minik8s/pkg/apiServer/etcdClient"
)

// CreatePV 创建PersistentVolume
func CreatePV(c *gin.Context) {
	pv := &apiObject.PersistentVolume{}
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

	key := config.EtcdPvPrefix + "/" + pvNamespace + "/" + pvName
	response, _ := etcdclient.EtcdStore.Get(key)
	if response != "" {
		log.ErrorLog("Create PersistentVolume: pv already exists" + response)
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "pv already exists"})
		return
	}
	log.InfoLog("CreatePv: " + key)

	// 将本地目录 /pvclient/:namespace/:name 挂载到服务器目录 /pvserver/:namespace/:name
	mountCmd := "mount -t nfs " + config.NFSServer + ":/pvserver/" + pvNamespace + "/" + pvName + " /pvclient/" + pvNamespace + "/" + pvName
	cmd := exec.Command("sh", "-c", mountCmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.ErrorLog("Create PersistentVolume: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.InfoLog("Create PersistentVolume: " + string(output))

	c.JSON(200, gin.H{"data": "Create PersistentVolume " + pvName})
}
