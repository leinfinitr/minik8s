package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/apiObject"
	etcdclient "minik8s/pkg/apiServer/etcdClient"
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

	// 检查 serverless 对象是否存储在 etcd
	key := config.EtcdServerlessPrefix + "/" + serverless.Name
	response, _ := etcdclient.EtcdStore.Get(key)
	if response != "" {
		log.ErrorLog("CreateServerless: " + serverless.Name + " already exists")
		c.JSON(400, gin.H{"error": "Serverless " + serverless.Name + " already exists"})
		return
	}

	// 根据 serverless 对象创建一个 pod 对象
	pod := conversion.ServerlessToPod(serverless)

	// TODO: 将函数文件存入 volume

	// 将 pod 对象存入 etcd
	podJson, err := json.Marshal(pod)
	err = etcdclient.EtcdStore.Put(key, string(podJson))
	if err != nil {
		log.ErrorLog("CreateServerless: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	log.InfoLog("CreateServerless: " + serverless.Name)
}

// DeleteServerless 删除Serverless环境
func DeleteServerless(c *gin.Context) {
	log.DebugLog("DeleteServerless")
	serverlessName := c.Param("serverlessName")
	if serverlessName == "" {
		log.ErrorLog("DeleteServerless: serverlessName is empty")
		c.JSON(400, gin.H{"error": "serverlessName is empty"})
		return
	}

	// 从 etcd 中删除 serverless 对象
	key := config.EtcdServerlessPrefix + "/" + serverlessName
	err := etcdclient.EtcdStore.Delete(key)
	if err != nil {
		log.ErrorLog("DeleteServerless: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	log.InfoLog("DeleteServerless: " + serverlessName)
}

// GetServerless 获取所有的Serverless Function
func GetServerless(c *gin.Context) {
	log.InfoLog("GetServerless")
	var podList []apiObject.Pod
	// 从 etcd 中获取所有属于 Serverless 的 Pod 对象
	response, err := etcdclient.EtcdStore.Get(config.EtcdServerlessPrefix)
	if err != nil {
		log.ErrorLog("GetServerless: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 将 json 字符串转换为 pod 对象列表
	err = json.Unmarshal([]byte(response), &podList)
	if err != nil {
		log.ErrorLog("GetServerless: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 将 pod 对象列表转换为 serverless 对象列表
	var serverlessList []apiObject.Serverless
	for _, pod := range podList {
		serverless := conversion.PodToServerless(pod)
		serverlessList = append(serverlessList, serverless)
	}
	c.JSON(200, gin.H{"data": serverlessList})
}

// RunServerlessFunction 运行Serverless Function
func RunServerlessFunction(c *gin.Context) {
	log.DebugLog("RunServerlessFunction")
	serverlessName := c.Param("serverlessName")
	if serverlessName == "" {
		log.ErrorLog("RunServerlessFunction: serverlessName or functionName is empty")
		c.JSON(400, gin.H{"error": "serverlessName or functionName is empty"})
		return
	}
	log.InfoLog("RunServerlessFunction: " + serverlessName)

	// 获取 serverless 对应的 pod 对象
	key := config.EtcdServerlessPrefix + "/" + serverlessName
	response, _ := etcdclient.EtcdStore.Get(key)
	if response == "" {
		log.ErrorLog("RunServerlessFunction: " + serverlessName + " not exists")
		c.JSON(400, gin.H{"error": "Serverless " + serverlessName + " not exists"})
		return
	}

	var pod apiObject.Pod
	err := json.Unmarshal([]byte(response), &pod)
	if err != nil {
		log.ErrorLog("RunServerlessFunction: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 转发给 apiServer 运行 pod
	go func() {
		url := config.APIServerURL() + config.PodsURI
		url = strings.Replace(url, config.NameSpaceReplace, pod.Metadata.Namespace, -1)
		res, err := httprequest.PostObjMsg(url, pod)
		if err != nil {
			log.ErrorLog("RunServerlessFunction: " + err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if res.StatusCode != 200 {
			log.ErrorLog("RunServerlessFunction: " + res.Status)
			c.JSON(500, gin.H{"error": res.Status})
			return
		} else {
			log.InfoLog("RunServerlessFunction success" + serverlessName)
			c.JSON(200, gin.H{"data": "success"})
		}
	}()
}

// UpdateServerlessFunction 更新Serverless Function
func UpdateServerlessFunction(c *gin.Context) {
	log.DebugLog("UpdateServerlessFunction")
	serverlessName := c.Param("serverlessName")
	if serverlessName == "" {
		log.ErrorLog("UpdateServerlessFunction: serverlessName or functionName is empty")
		c.JSON(400, gin.H{"error": "serverlessName or functionName is empty"})
		return
	}
	log.InfoLog("UpdateServerlessFunction: " + serverlessName)

	// 得到 PUT 请求中的文件名
	var fileName string
	err := c.ShouldBindJSON(&fileName)
	if err != nil {
		log.ErrorLog("UpdateServerlessFunction: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// TODO: 更新 volume 中的文件

	c.JSON(200, gin.H{"data": "success"})
}
