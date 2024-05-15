package handler

import (
	"github.com/gin-gonic/gin"
	"minik8s/pkg/apiObject"
	"minik8s/tools/log"
)

// CreateServerless 创建Serverless环境
func CreateServerless(c *gin.Context) {
	log.DebugLog("CreateServerless")
	var serverless = apiObject.Serverless{}
	err := c.ShouldBindJSON(&serverless)
	if err != nil {
		log.ErrorLog("CreateServerless: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
	}
}

// GetServerless 获取所有的Serverless Function
//func GetServerless(c *gin.Context) {
//	name := c.Param("name")
//	namespace := c.Param("namespace")
//	if namespace == "" {
//		namespace = "default"
//	} else if name == "" {
//		log.ErrorLog("GetPod: name is empty")
//		c.JSON(400, gin.H{"error": "name is empty"})
//		return
//	}
//	log.InfoLog("GetPod: " + namespace + "/" + name)
//	key := config.EtcdPodPrefix + "/" + namespace + "/" + name
//	res, err := etcdclient.EtcdStore.Get(key)
//	if err != nil {
//		log.ErrorLog("GetPod: " + err.Error())
//		c.JSON(500, gin.H{"error": err.Error()})
//		return
//	}
//	resJson, err := json.Marshal(res)
//	if err != nil {
//		log.ErrorLog("GetPod: " + err.Error())
//		c.JSON(500, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(200, gin.H{"data": resJson})
//}

// CreateServerless 创建Serverless Function
//func CreateServerless(c *gin.Context) {
//	image := c.PostForm("image")
//	name := c.Param("name")
//	functionList := c.PostFormArray("function")
//
//	if namespace == "" {
//		namespace = "default"
//	} else if name == "" {
//		log.ErrorLog("CreatePod: name is empty")
//		c.JSON(400, gin.H{"error": "name is empty"})
//		return
//	}
//	log.InfoLog("CreatePod: " + namespace + "/" + name)
//	key := config.EtcdPodPrefix + "/" + namespace + "/" + name
//	_, err := etcdclient.EtcdStore.Get(key)
//	if err != nil {
//		log.ErrorLog("CreatePod: " + err.Error())
//		c.JSON(500, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(200, gin.H{"data": "success"})
//}
