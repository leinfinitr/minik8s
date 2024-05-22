package handlers

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/pkg/entity"
	httprequest "minik8s/tools/httpRequest"
	"minik8s/tools/log"

	etcdclient "minik8s/pkg/apiServer/etcdClient"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TODO: UpdateProxy 更新Proxy的状态
func UpdateProxyStatus(c *gin.Context) {
	// var node apiObject.Node
	// err := c.ShouldBindJSON(&node)
	// if err != nil {
	// 	log.ErrorLog("UpdateNode error: " + err.Error())
	// }
	// name := c.Param("name")
	// nodes[name] = node

	// log.InfoLog("UpdateNode: " + name)
	// c.JSON(config.HttpSuccessCode, "")
}

// GetServices 获取所有Service
func GetServices(c *gin.Context) {
	namespace := c.Param("namespace")

	println("GetServices: " + namespace)
}

// GetService 获取指定Service
func GetService(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")

	println("GetService: " + namespace + "/" + name)
}

// GetService 获取指定Service
func PutService(c *gin.Context) {
	var serviceEvent entity.ServiceEvent

	service := &apiObject.Service{}
	err := c.ShouldBindJSON(service)
	if err != nil {
		log.ErrorLog("PutService error: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	newServiceName := service.Metadata.Name
	newServiceNamespace := service.Metadata.Namespace
	if newServiceNamespace == "" || newServiceName == "" {
		log.ErrorLog("PutService error: namespace or name is empty")
		c.JSON(400, gin.H{"error": "namespace or name is empty"})
		return
	}
	key := config.EtcdServicePrefix + "/" + newServiceNamespace + "/" + newServiceName
	response, _ := etcdclient.EtcdStore.Get(key)
	if response != "" {
		serviceEvent.Action = "UpdateService"
	} else {
		serviceEvent.Action = "CreateService"
	}

	service.Metadata.UUID = uuid.New().String()
	serviceEvent.Service = *service
	serviceEvent.Endpoints = *Selector(service)

	resJson, err := json.Marshal(serviceEvent.Service)
	if err != nil {
		log.WarnLog("GetNodes: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}
	etcdclient.EtcdStore.Put(key, string(resJson))
	log.InfoLog("PutService: " + newServiceNamespace + "/" + newServiceName)

	// 获取所有的Node信息
	res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdNodePrefix)
	if err != nil {
		log.WarnLog("GetNodes: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}

	var nodes []apiObject.Node
	for _, v := range res {
		var node apiObject.Node
		err = json.Unmarshal([]byte(v), &node)
		if err != nil {
			log.WarnLog("GetNodes: " + err.Error())
			c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
			return
		}
		nodes = append(nodes, node)
	}

	// 向所有的Node发送serviceEvent
	for _, node := range nodes {
		url := "http://" + node.Status.Addresses[0].Address + ":" + fmt.Sprint(config.KubeproxyAPIPort) + "/" + config.ServiceURI
		if serviceEvent.Action == "CreateService" {
			res, err := httprequest.PutObjMsg(url, serviceEvent)
			if err != nil || res.StatusCode != 200 {
				log.ErrorLog("PutService: " + err.Error())
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		} else {
			res, err := httprequest.PostObjMsg(url, serviceEvent)
			if err != nil || res.StatusCode != config.HttpSuccessCode {
				log.ErrorLog("PostService: " + err.Error())
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}
	}

	log.InfoLog("Put Service Successfully")
	c.JSON(config.HttpSuccessCode, "Service add successfully")
	return
}

// GetServiceStatus 获取指定Service的状态
func GetServiceStatus(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")

	println("GetServiceStatus: " + namespace + "/" + name)
}

// 从etcd中获取所有和service相关的pod
func Selector(service *apiObject.Service) *[]apiObject.Endpoint {
	var endpoints []apiObject.Endpoint
	selector := service.Spec.Selector
	key := config.EtcdPodPrefix + "/" + service.Metadata.Namespace
	res, err := etcdclient.EtcdStore.PrefixGet(key)
	if err != nil {
		log.ErrorLog("GetPods: " + err.Error())
		return nil
	}
	for _, str := range res {
		pod := apiObject.Pod{}
		json.Unmarshal([]byte(str), &pod)
		flag := true
		for k, v := range selector {
			if pod.Metadata.Labels[k] != v {
				flag = false
				break
			}
		}
		if flag {
			endpoint := apiObject.Endpoint{
				PodUUID: pod.APIVersion,
				IP:      pod.Status.PodIP,
			}
			endpoints = append(endpoints, endpoint)
		}
	}

	return &endpoints
}
