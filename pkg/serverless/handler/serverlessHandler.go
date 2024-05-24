package handler

import (
	"encoding/json"
	"minik8s/pkg/apiObject"
	etcdclient "minik8s/pkg/apiServer/etcdClient"
	"minik8s/pkg/config"
	"minik8s/pkg/serverless/scale"
	"minik8s/tools/conversion"
	"minik8s/tools/log"

	"github.com/gin-gonic/gin"
)

// CreateServerless 创建Serverless环境
func CreateServerless(c *gin.Context) {
	log.DebugLog("CreateServerless")
	var serverless = apiObject.Serverless{}
	err := c.ShouldBindJSON(&serverless)
	if err != nil {
		log.ErrorLog("CreateServerless: " + err.Error())
		c.JSON(400, err.Error())
	}

	// 检查 serverless 对象是否存储在 etcd
	key := config.EtcdServerlessPrefix + "/" + serverless.Name
	response, _ := etcdclient.EtcdStore.Get(key)
	if response != "" {
		log.ErrorLog("CreateServerless: " + serverless.Name + " already exists")
		c.JSON(400, "Serverless "+serverless.Name+" already exists")
		return
	}

	// 根据 serverless 对象创建一个 pod 对象
	pod := conversion.ServerlessToPod(serverless)

	// TODO: 将函数文件存入 volume

	// 将 pod 对象存入 etcd
	podJson, err := json.Marshal(pod)
	if err != nil {
		log.ErrorLog("CreateServerless: " + err.Error())
		c.JSON(500, err.Error())
		return
	}
	err = etcdclient.EtcdStore.Put(key, string(podJson))
	if err != nil {
		log.ErrorLog("CreateServerless: " + err.Error())
		c.JSON(500, err.Error())
		return
	}

	log.InfoLog("CreateServerless: " + serverless.Name)
	c.JSON(200, "Create serverless "+serverless.Name+" success")
}

// GetServerless 获取所有的Serverless Function
func GetServerless(c *gin.Context) {
	log.InfoLog("GetServerless")
	var podList []apiObject.Pod
	// 从 etcd 中获取所有属于 Serverless 的 Pod 对象
	response, err := etcdclient.EtcdStore.PrefixGet(config.EtcdServerlessPrefix)
	if err != nil {
		log.ErrorLog("GetServerless: " + err.Error())
		c.JSON(500, err.Error())
		return
	}
	// 遍历 response，依次将 json 字符串转换为 pod 对象
	for _, podJson := range response {
		var pod apiObject.Pod
		err = json.Unmarshal([]byte(podJson), &pod)
		if err != nil {
			log.ErrorLog("GetServerless: " + err.Error())
			c.JSON(500, err.Error())
			return
		}
		podList = append(podList, pod)
	}
	// 将 pod 对象列表转换为 serverless 对象列表
	var serverlessList []apiObject.Serverless
	for _, pod := range podList {
		serverless := conversion.PodToServerless(pod)
		serverlessList = append(serverlessList, serverless)
	}
	c.JSON(200, serverlessList)
}

// DeleteServerless 删除Serverless环境
func DeleteServerless(c *gin.Context) {
	log.DebugLog("DeleteServerless")
	serverlessName := c.Param("name")
	if serverlessName == "" {
		log.ErrorLog("DeleteServerless: serverlessName is empty")
		c.JSON(400, "ServerlessName is empty")
		return
	}

	// 从 etcd 中删除 serverless 对象
	key := config.EtcdServerlessPrefix + "/" + serverlessName
	err := etcdclient.EtcdStore.Delete(key)
	if err != nil {
		log.ErrorLog("DeleteServerless: " + err.Error())
		c.JSON(500, err.Error())
		return
	}

	log.InfoLog("DeleteServerless: " + serverlessName)
	c.JSON(200, "Delete serverless "+serverlessName+" success")
}

// UpdateServerlessFunction 更新Serverless Function
func UpdateServerlessFunction(c *gin.Context) {
	log.DebugLog("UpdateServerlessFunction")
	serverlessName := c.Param("name")
	if serverlessName == "" {
		log.ErrorLog("UpdateServerlessFunction: serverlessName is empty")
		c.JSON(400, "ServerlessName is empty")
		return
	}
	log.InfoLog("UpdateServerlessFunction: " + serverlessName)

	// 得到 PUT 请求中的文件名
	var fileName string
	err := c.ShouldBindJSON(&fileName)
	if err != nil {
		log.ErrorLog("UpdateServerlessFunction: " + err.Error())
		c.JSON(400, err.Error())
		return
	}

	// TODO: 更新 volume 中的文件

	c.JSON(200, "Update serverless "+serverlessName+" success")
}

// RunServerlessFunction 运行Serverless Function
func RunServerlessFunction(c *gin.Context) {
	log.DebugLog("RunServerlessFunction")
	serverlessName := c.Param("name")
	param := c.Param("param")
	if serverlessName == "" {
		log.ErrorLog("RunServerlessFunction: serverlessName or functionName is empty")
		c.JSON(400, "ServerlessName or functionName is empty")
		return
	}
	log.InfoLog("RunServerlessFunction: " + serverlessName)

	result := RunFunction(serverlessName, param)
	c.JSON(200, result)
}

// RunFunction 运行Serverless Function
func RunFunction(name string, param string) string {
	// 获取 serverless 对应的 pod 对象
	key := config.EtcdServerlessPrefix + "/" + name
	response, _ := etcdclient.EtcdStore.Get(key)
	if response == "" {
		log.ErrorLog("RunFunction: " + name + " not found")
		return ""
	}
	var pod apiObject.Pod
	err := json.Unmarshal([]byte(response), &pod)
	if err != nil {
		log.ErrorLog("RunFunction: " + err.Error())
		return ""
	}

	// 运行请求数加一
	scale.ScaleManager.IncreaseRequestNum(name)

	// 交由 scale 进行处理
	result := scale.ScaleManager.RunFunction(pod, param)

	// 运行请求数减一
	scale.ScaleManager.DecreaseRequestNum(name)

	return result
}
