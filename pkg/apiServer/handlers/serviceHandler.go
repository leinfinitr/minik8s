package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/pkg/entity"
	"minik8s/tools/log"

	etcdclient "minik8s/pkg/apiServer/etcdClient"
	httprequest "minik8s/tools/httpRequest"
)

func RegisterProxy(c *gin.Context) {
	// 某个proxy初次注册，检查是否已经有service存在，如果有则将service发送给proxy
	res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdServicePrefix)
	if err != nil {
		log.ErrorLog("RegisterProxy: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}

	var services []apiObject.Service
	for _, v := range res {
		var service apiObject.Service
		err = json.Unmarshal([]byte(v), &service)
		if err != nil {
			log.ErrorLog("RegisterProxy: " + err.Error())
			c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
			return
		}
		services = append(services, service)
	}

	// 如果kubeproxy所在的node没有注册，则返回错误
	var node apiObject.Node
	res, err = etcdclient.EtcdStore.PrefixGet(config.EtcdNodePrefix)
	if err != nil {
		log.ErrorLog("RegisterProxy: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}
	for _, v := range res {
		err = json.Unmarshal([]byte(v), &node)
		if err != nil {
			log.ErrorLog("RegisterProxy: " + err.Error())
			c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
			return
		}
		if node.Status.Addresses[0].Address == c.ClientIP() {
			break
		}
	}

	if node.Status.Addresses[0].Address != c.ClientIP() {
		log.ErrorLog("RegisterProxy: node not exists, IP: " + c.ClientIP() + " node: " + node.Metadata.Name)
		c.JSON(config.HttpErrorCode, gin.H{"error": "node not exists"})
		return
	}

	// 向proxy发送serviceEvent
	for _, service := range services {
		var serviceEvent entity.ServiceEvent
		serviceEvent.Action = entity.UpdateEvent
		serviceEvent.Service = service
		serviceEvent.Endpoints = *Selector(&service)

		// 向proxy发送serviceEvent
		url := "http://" + c.ClientIP() + ":" + fmt.Sprint(config.KubeproxyAPIPort) + config.ServiceURI
		url = strings.Replace(url, config.NameSpaceReplace, service.Metadata.Namespace, -1)
		url = strings.Replace(url, config.NameReplace, service.Metadata.Name, -1)
		res, err := httprequest.PostObjMsg(url, serviceEvent)
		if err != nil || res.StatusCode != config.HttpSuccessCode {
			log.ErrorLog("RegisterProxy: " + err.Error())
			c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
			return
		}

	}

	log.InfoLog("RegisterProxy Successfully")
	c.JSON(config.HttpSuccessCode, "RegisterProxy Successfully")

}

// GetServices 获取所有Service
func GetServices(c *gin.Context) {

	res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdServicePrefix)
	if err != nil {
		log.ErrorLog("GetServices: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}

	var services []apiObject.Service
	for _, v := range res {
		var service apiObject.Service
		err = json.Unmarshal([]byte(v), &service)
		if err != nil {
			log.ErrorLog("RegisterProxy: " + err.Error())
			c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
			return
		}
		services = append(services, service)
	}

	c.JSON(config.HttpSuccessCode, services)

}

// GetService 获取指定Service
func GetService(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")

	println("GetService: " + namespace + "/" + name)
}

func DeleteService(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	key := config.EtcdServicePrefix + "/" + namespace + "/" + name
	response, err := etcdclient.EtcdStore.Get(key)
	if response == "" || err != nil {
		log.ErrorLog("DeleteService error: service " + namespace + "/" + name + " not exists")
		c.JSON(400, gin.H{"error": "service not exists"})
		return
	}
	var service apiObject.Service
	err = json.Unmarshal([]byte(response), &service)
	if err != nil {
		log.ErrorLog("DeleteService error: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 把deleteEvent发送给所有的Node
	var serviceEvent entity.ServiceEvent
	serviceEvent.Action = entity.DeleteEvent
	serviceEvent.Service = service
	serviceEvent.Endpoints = *Selector(&service)

	// 获取所有的Node信息
	res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdNodePrefix)
	if err != nil {
		log.WarnLog("GetNodes: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}
	for _, v := range res {
		var node apiObject.Node
		err = json.Unmarshal([]byte(v), &node)
		if err != nil {
			log.WarnLog("GetNodes: " + err.Error())
			c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
			return
		}
		url := "http://" + node.Status.Addresses[0].Address + ":" + fmt.Sprint(config.KubeproxyAPIPort) + config.ServiceURI
		url = strings.Replace(url, config.NameSpaceReplace, namespace, -1)
		url = strings.Replace(url, config.NameReplace, name, -1)
		res, err := httprequest.PostObjMsg(url, serviceEvent)
		if err != nil || res.StatusCode != config.HttpSuccessCode {
			log.ErrorLog("DeleteService: " + err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	etcdclient.EtcdStore.Delete(key)

	// 删除service2endpoint
	key = config.EtcdService2EndpointPrefix + "/" + namespace + "/" + name
	etcdclient.EtcdStore.Delete(key)
}

// PutService 获取指定Service
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
		serviceEvent.Action = entity.UpdateEvent
	} else {
		serviceEvent.Action = entity.CreateEvent
	}

	service.Metadata.UUID = uuid.New().String()
	service.Spec.ClusterIP = AllocClusterIP()
	log.InfoLog("AllocClusterIP: " + service.Spec.ClusterIP)
	serviceEvent.Service = *service
	serviceEvent.Endpoints = *Selector(service)

	resJson, err := json.Marshal(serviceEvent.Service)
	if err != nil {
		log.WarnLog("GetNodes: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(key, string(resJson))
	if err != nil {
		log.WarnLog("GetNodes: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}
	log.InfoLog("PutService: " + newServiceNamespace + "/" + newServiceName)

	// 将service存入etcd
	service2Endpoint, err := json.Marshal(serviceEvent)
	if err != nil {
		log.WarnLog("GetNodes: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}
	key = config.EtcdService2EndpointPrefix + "/" + newServiceNamespace + "/" + newServiceName
	if err = etcdclient.EtcdStore.Put(key, string(service2Endpoint)); err != nil {
		log.WarnLog("GetNodes: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}

	nodes := GetALLNodes()

	// 向所有的Node发送serviceEvent
	for _, node := range nodes {
		url := "http://" + node.Status.Addresses[0].Address + ":" + fmt.Sprint(config.KubeproxyAPIPort) + config.ServiceURI
		url = strings.Replace(url, config.NameSpaceReplace, newServiceNamespace, -1)
		url = strings.Replace(url, config.NameReplace, newServiceName, -1)
		if serviceEvent.Action == entity.CreateEvent {
			res, err := httprequest.PutObjMsg(url, serviceEvent)
			if err != nil || res.StatusCode != 200 {
				log.ErrorLog("PutService: errorUrl: " + url + "error: " + err.Error())
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

}

// GetServiceStatus 获取指定Service的状态
func GetServiceStatus(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")

	println("GetServiceStatus: " + namespace + "/" + name)
}

// Selector 从etcd中获取所有和service相关的pod
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

func AllocClusterIP() string {
	IP := ""
	var service apiObject.Service
	isused := false

	for {
		i := 0
		for i = 0; i < 4; i++ {
			IP += fmt.Sprint(rand.Intn(255))
			if i != 3 {
				IP += "."
			}
		}

		// 比对所有的service的clusterIP，如果有重复则重新生成
		res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdServicePrefix)
		if err != nil {
			log.ErrorLog("AllocClusterIP: " + err.Error())
			return ""
		}

		for _, v := range res {
			json.Unmarshal([]byte(v), &service)
			if service.Spec.ClusterIP == IP {
				isused = true
				break
			}
		}

		if !isused {
			break
		}
	}

	return IP
}
