package kubeproxy

import (
	"fmt"
	"time"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/pkg/entity"
	"minik8s/pkg/kubeproxy/iptableManager"
	"minik8s/tools/host"
	"minik8s/tools/log"
	"minik8s/tools/netRequest"

	"github.com/gin-gonic/gin"
)

type Kubeproxy struct {
	// 从UUID到service的map
	serviceUUIDMap  map[string]apiObject.Service
	proxyAPIRouter  *gin.Engine
	apiServerConfig config.APIServerConfig
	iptableManager  iptableManager.IptableManager
}

var kubeproxy *Kubeproxy

func GetKubeproxy() *Kubeproxy {
	if kubeproxy == nil {
		kubeproxy = &Kubeproxy{
			serviceUUIDMap: make(map[string]apiObject.Service),
			// serviceEvents:   make(chan *entity.ServiceEvent),
			proxyAPIRouter:  gin.Default(),
			apiServerConfig: *config.NewAPIServerConfig(),
			// iptableManager: GetIptableManager(),
		}
	}
	return kubeproxy
}

func (k *Kubeproxy) createService(c *gin.Context) {
	var serviceEvent *entity.ServiceEvent
	err := c.ShouldBindJSON(&serviceEvent)
	if err != nil {
		log.ErrorLog("CreateService error: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}

	service := serviceEvent.Service

	// 保存service
	if _, ok := k.serviceUUIDMap[service.Metadata.UUID]; ok {
		log.ErrorLog("CreateService error: service already exists")
		c.JSON(config.HttpErrorCode, gin.H{"error": "service already exists"})
		return
	}
	k.serviceUUIDMap[service.Metadata.UUID] = service

	// 调用iptableManager创建service
	err = k.iptableManager.CreateService(serviceEvent)
	if err != nil {
		log.ErrorLog("CreateService error: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(config.HttpSuccessCode, "Service created successfully")

}

func (k *Kubeproxy) updateService(c *gin.Context) {
	var serviceEvent *entity.ServiceEvent
	err := c.ShouldBindJSON(&serviceEvent)
	if err != nil {
		log.ErrorLog("UpdateService error: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}

	service := serviceEvent.Service

	// 更新service
	if _, ok := k.serviceUUIDMap[service.Metadata.UUID]; !ok {
		log.ErrorLog("UpdateService error: service not exists")
		c.JSON(config.HttpErrorCode, gin.H{"error": "service not exists"})
		return
	}
	k.serviceUUIDMap[service.Metadata.UUID] = service

	// 调用iptableManager更新service
	err = k.iptableManager.UpdateService(serviceEvent)
	if err != nil {
		log.ErrorLog("UpdateService error: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(config.HttpSuccessCode, "Service updated successfully")

}

func (k *Kubeproxy) deleteService(c *gin.Context) {
	var serviceEvent *entity.ServiceEvent
	err := c.ShouldBindJSON(&serviceEvent)
	if err != nil {
		log.ErrorLog("DeleteService error: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}

	service := serviceEvent.Service

	// 删除service
	if _, ok := k.serviceUUIDMap[service.Metadata.UUID]; !ok {
		log.ErrorLog("DeleteService error: service not exists")
		c.JSON(config.HttpErrorCode, gin.H{"error": "service not exists"})
		return
	}

	// 调用iptableManager删除service
	err = k.iptableManager.DeleteService(serviceEvent)
	if err != nil {
		log.ErrorLog("DeleteService error: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}
	delete(k.serviceUUIDMap, service.Metadata.UUID)

	c.JSON(config.HttpSuccessCode, "Service deleted successfully")

}

func (k *Kubeproxy) registerKubeproxyAPI() {
	// 创建一个新的service
	k.proxyAPIRouter.PUT(config.ServiceURI, k.createService)
	// 更新指定的service
	k.proxyAPIRouter.POST(config.ServiceURI, k.updateService)
	// 删除制定的service
	k.proxyAPIRouter.DELETE(config.ServiceURI, k.deleteService)
}

func (k *Kubeproxy) Run() {

	// 主线程，用于接受并转发来自与apiServer通信端口的请求
	k.registerKubeproxyAPI()
	kubeproxyIP, _ := host.GetHostIP()
	_ = k.proxyAPIRouter.Run(kubeproxyIP + ":" + fmt.Sprint(config.KubeproxyAPIPort))

	// 创建用于发送心跳的线程
	go k.heartbeat()
}

func (k *Kubeproxy) heartbeat() {
	log.InfoLog("[Kubeproxy] ProxyStatus")

	// 每间隔60s发送一次心跳
	for {
		HostName, _ := host.GetHostname()
		hostIP, _ := host.GetHostIP()
		url := k.apiServerConfig.APIServerURL() + config.ProxyStatusURI + "/" + HostName + "/status"

		statusCode, _, _ := netRequest.PutRequestByTarget(url, hostIP)

		if statusCode != config.HttpSuccessCode {
			log.ErrorLog("kubeproxy heartbeat failed")
		} else {
			log.DebugLog("kubeproxy heartbeat success")
		}

		time.Sleep(60 * time.Second)
	}

}
