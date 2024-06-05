package kubeproxy

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
			proxyAPIRouter:  gin.New(),
			apiServerConfig: *config.NewAPIServerConfig(),
			iptableManager:  iptableManager.GetIptableManager(),
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
		log.WarnLog("UpdateService error: service not exists")
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

func (k *Kubeproxy) createDNS(c *gin.Context) {
	// 把该域名和IP的映射增加到本地的/etc/hosts文件中
	var dns apiObject.Dns
	if err := c.ShouldBindJSON(&dns); err != nil {
		log.ErrorLog("CreateDNS error: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}

	filePath := config.DNS_PATH
	ip := dns.NginxIP
	domain := dns.Spec.Host

	// 读取文件
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 读取文件的每一行
	scanner := bufio.NewScanner(file)
	var lines []string
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, domain) {
			// 如果找到了域名，删除旧域名
			continue
		} else if strings.Contains(line, ip) {
			// 如果找到了 IP，追加域名
			line = line + " " + domain
			found = true
		}
		lines = append(lines, line)
	}

	// 如果没有找到 IP，添加新的映射
	if !found {
		lines = append(lines, fmt.Sprintf("%s %s", ip, domain))
	}

	// 写回文件
	file, err = os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()

}

func (k *Kubeproxy) deleteDNS(c *gin.Context) {
	// 把该域名和IP的映射从本地的/etc/hosts文件中删除
	var dns apiObject.Dns
	if err := c.ShouldBindJSON(&dns); err != nil {
		log.ErrorLog("DeleteDNS error: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}

	filePath := config.DNS_PATH
	ip := dns.NginxIP
	domain := dns.Spec.Host

	// 读取文件
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 读取文件的每一行
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		// 如果找到了 IP 和域名
		if strings.Contains(line, ip) && strings.Contains(line, domain) {
			// 删除域名
			line = strings.Replace(line, domain, "", -1)
			// 如果该行只有这个 IP，则不添加到 lines 中
			if strings.TrimSpace(line) == ip {
				continue
			}
		}
		lines = append(lines, line)
	}

	// 写回文件
	file, err = os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()
	c.JSON(config.HttpSuccessCode, "DNS deleted successfully")

}

func (k *Kubeproxy) registerKubeproxyAPI() {
	// 创建一个新的service
	k.proxyAPIRouter.PUT(config.ServiceURI, k.createService)
	// 更新指定的service
	k.proxyAPIRouter.POST(config.ServiceURI, k.updateService)
	// 删除制定的service
	k.proxyAPIRouter.DELETE(config.ServiceURI, k.deleteService)

	// 与DNS相关的部分
	// 创建一个新的DNS
	k.proxyAPIRouter.PUT(config.DNSURI, k.createDNS)
	// 删除指定的DNS
	k.proxyAPIRouter.DELETE(config.DNSURI, k.deleteDNS)
}

func (k *Kubeproxy) Run() {

	// 主线程，用于接受并转发来自与apiServer通信端口的请求
	k.registerKubeproxyAPI()
	kubeproxyIP, _ := host.GetHostIP()
	// 在proxy刚启动时，向apiServer注册自己
	go k.registerProxy()

	_ = k.proxyAPIRouter.Run(kubeproxyIP + ":" + fmt.Sprint(config.KubeproxyAPIPort))
}

func (k *Kubeproxy) registerProxy() {
	log.InfoLog("[Kubeproxy] Register to apiServer")

	// 一直尝试直到注册成功
	for {
		url := k.apiServerConfig.APIServerURL() + config.ProxyStatusURI

		statusCode, _, _ := netRequest.PostRequestByTarget(url, nil)

		if statusCode != config.HttpSuccessCode {
			log.ErrorLog("kubeproxy register failed")
		} else {
			log.DebugLog("kubeproxy register success")
			return
		}

		time.Sleep(15 * time.Second)
	}

}
