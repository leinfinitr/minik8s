package handler

import (
	"github.com/gin-gonic/gin"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/conversion"
	"minik8s/tools/httpRequest"
	"minik8s/tools/log"
	"strings"
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

	// 根据 serverless 对象创建一个 pod 对象
	pod := conversion.ServerlessToPod(serverless)

	// 转发给 apiServer 创建 pod
	url := config.APIServerURL() + config.PodsURI
	url = strings.Replace(url, config.NameSpaceReplace, pod.Metadata.Namespace, -1)
	_, err = httprequest.PostObjMsg(config.APIServerURL()+config.PodsURI, pod)
	if err != nil {
		log.ErrorLog("CreateServerless: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	log.DebugLog("CreateServerless success" + serverless.Name)

	// 将函数文件存入 volume
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
