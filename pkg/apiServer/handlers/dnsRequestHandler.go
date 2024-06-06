package handlers

import (
	"github.com/gin-gonic/gin"

	"minik8s/pkg/config"
	"minik8s/tools/log"

	etcdclient "minik8s/pkg/apiServer/etcdClient"
)

func GetGlobalDnsRequests(c *gin.Context) {
	//log.InfoLog("GetGlobalDnsRequests")
	//res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdDnsRequestPrefix)
	//if err != nil {
	//	log.ErrorLog("GetGlobalDnsRequests: " + err.Error())
	//	c.JSON(500, gin.H{"error": err.Error()})
	//	return
	//}
	//c.JSON(200, resJson)
}

func DeleteDnsRequest(c *gin.Context) {
	log.InfoLog("DeleteDnsRequest")
	name := c.Param("name")
	namespace := c.Param("namespace")
	if name == "" {
		log.ErrorLog("DeleteDnsRequest name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	if namespace == "" {
		namespace = "default"
	}
	log.InfoLog("DeleteDnsRequest: " + namespace + "/" + name)

	key := config.EtcdDnsRequestPrefix + "/" + namespace + "/" + name
	resJson, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("DeleteDnsRequest: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(resJson) == 0 {
		log.ErrorLog("DeleteDnsRequest: not found")
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	err = etcdclient.EtcdStore.Delete(key)
	if err != nil {
		log.ErrorLog("DeleteDnsRequest: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": resJson})
}
