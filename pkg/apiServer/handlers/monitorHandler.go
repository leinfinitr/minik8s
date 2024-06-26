package handlers

import (
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v3"

	"github.com/gin-gonic/gin"

	"minik8s/pkg/apiObject"
	"minik8s/tools/log"

	Config "minik8s/pkg/config"
)

func GetPrometheusConfig() *apiObject.PrometheusConfig {
	// 将nodeIp和端口信息追加到prometheus配置文件中
	// 1. 读取prometheus配置文件
	data, err := os.ReadFile(apiObject.PrometheusConfigPath)
	if err != nil {
		log.ErrorLog("Failed to read Prometheus config file: " + err.Error())
		return nil
	}
	// 2. 追加nodeIp和端口信息
	var config apiObject.PrometheusConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.ErrorLog("Failed to parse Prometheus config file: " + err.Error())
		return nil
	}

	return &config
}

func PutPrometheusConfig(config *apiObject.PrometheusConfig) error {
	// 保存并写入配置文件
	newData, err := yaml.Marshal(&config)
	if err != nil {
		log.ErrorLog("Failed to marshal Prometheus config file: " + err.Error())
		return err
	}

	err = os.WriteFile(apiObject.PrometheusConfigPath, newData, 0644)
	if err != nil {
		log.ErrorLog("Failed to write Prometheus config file: " + err.Error())
		return err
	}

	return nil

}

// RegisterNodeMonitor 在首次注册时，在本地节点的prometheus配置文件中追加一个监控项
func RegisterNodeMonitor(c *gin.Context) {
	log.InfoLog("Start RegisterMonitor")

	// 1.从请求中获取Node的值
	var node apiObject.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 2. 获取node的IP
	nodeIP := node.Status.Addresses[0].Address

	config := GetPrometheusConfig()
	if config == nil {
		c.JSON(500, gin.H{"error": "Failed to get Prometheus config"})
		return
	}

	// 3. 按需添加新的job
	newTarget := nodeIP + ":" + fmt.Sprint(apiObject.NodeExporterPort)
	found := false
	for _, scrapeConfig := range config.ScrapeConfigs {
		if scrapeConfig.StaticConfigs[0].Targets[0] == newTarget {
			found = true
			break
		}
	}
	if !found {
		// 不存在则创建新的job
		NewStaticConfig := apiObject.StaticConfig{
			Targets: []string{newTarget},
			Labels:  map[string]string{"instance": node.Metadata.Name},
		}

		newScrapeConfig := apiObject.ScrapeConfig{
			JobName:       "node-exporter-" + node.Metadata.Name,
			StaticConfigs: []apiObject.StaticConfig{NewStaticConfig},
		}

		config.ScrapeConfigs = append(config.ScrapeConfigs, newScrapeConfig)
	} else {
		log.DebugLog("Node already registered")
		c.JSON(200, gin.H{"message": "Node already registered"})
		return
	}

	// 4. 保存并写入配置文件
	err := PutPrometheusConfig(config)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to write Prometheus config"})
		return
	}

	// 5. 热加载prometheus
	// curl -X POST http://192.168.1.7:9090/-/reload
	err = exec.Command("curl", "-X", "POST", "http://192.168.1.7:9090/-/reload").Run()
	if err != nil {
		log.ErrorLog("Failed to reload Prometheus: " + err.Error())
		c.JSON(500, gin.H{"error": "Failed to reload Prometheus"})
		return
	}

	c.JSON(200, gin.H{"message": "monitor registered successfully"})
}

func DeleteNodeMonitor(c *gin.Context) {
	// 1.从请求中获取Node的值
	var node apiObject.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 2. 获取node的IP
	nodeIP := node.Status.Addresses[0].Address
	if nodeIP == Config.APIServerLocalAddress {
		// 如果是apiServer节点，则不需要监控
		c.JSON(200, gin.H{"message": "Node is apiServer"})
		return
	}

	config := GetPrometheusConfig()
	if config == nil {
		c.JSON(500, gin.H{"error": "Failed to get Prometheus config"})
		return
	}

	// 3. 按需删除job
	newTarget := nodeIP + ":" + fmt.Sprint(apiObject.NodeExporterPort)
	found := false
	for i, scrapeConfig := range config.ScrapeConfigs {
		if scrapeConfig.StaticConfigs[0].Targets[0] == newTarget {
			config.ScrapeConfigs = append(config.ScrapeConfigs[:i], config.ScrapeConfigs[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		log.WarnLog("Node not registered")
		c.JSON(200, gin.H{"message": "Node not registered"})
		return
	}

	// 4. 保存并写入配置文件
	err := PutPrometheusConfig(config)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to write Prometheus config"})
		return
	}

	// 5. 热加载prometheus
	// curl -X POST http://192.168.1.7:9090/-/reload
	err = exec.Command("curl", "-X", "POST", "http://192.168.1.7:9090/-/reload").Run()
	if err != nil {
		log.ErrorLog("Failed to reload Prometheus: " + err.Error())
		c.JSON(500, gin.H{"error": "Failed to reload Prometheus"})
		return
	}

	c.JSON(200, gin.H{"message": "Node unregistered successfully"})

}

func RegisterPodMonitor(c *gin.Context) {
	var monitorPod apiObject.MonitorPod
	if err := c.ShouldBindJSON(&monitorPod); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 1. 读取prometheus配置文件
	config := GetPrometheusConfig()
	if config == nil {
		c.JSON(500, gin.H{"error": "Failed to get Prometheus config"})
		return
	}

	// 2. 在制定的job后面追加新的staticConfig，默认“pods”job是用来监控所有pod的
	var podScrapeConfig *apiObject.ScrapeConfig
	for _, scrapeConfig := range config.ScrapeConfigs {
		if scrapeConfig.JobName == "pods" {
			podScrapeConfig = &scrapeConfig
			break
		}
	}

	newStaticConfig := apiObject.StaticConfig{
		Targets: []string{},
		Labels:  map[string]string{"instance": monitorPod.PodName},
	}

	for _, uri := range monitorPod.MonitorUris {
		newStaticConfig.Targets = append(newStaticConfig.Targets, uri)
	}

	if podScrapeConfig == nil {
		// 不存在则创建新的job
		podScrapeConfig = &apiObject.ScrapeConfig{
			JobName:       "pods",
			StaticConfigs: []apiObject.StaticConfig{},
		}

		podScrapeConfig.StaticConfigs = append(podScrapeConfig.StaticConfigs, newStaticConfig)
		config.ScrapeConfigs = append(config.ScrapeConfigs, *podScrapeConfig)
	} else {
		podScrapeConfig.StaticConfigs = append(podScrapeConfig.StaticConfigs, newStaticConfig)
	}

	// 3. 保存并写入配置文件
	err := PutPrometheusConfig(config)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to write Prometheus config"})
		return
	}

	// 4. 热加载prometheus
	// curl -X POST http://192.168.1.7:9090/-/reload
	err = exec.Command("curl", "-X", "POST", "http://192.168.1.7:9090/-/reload").Run()
	if err != nil {
		log.ErrorLog("Failed to reload Prometheus: " + err.Error())
		c.JSON(500, gin.H{"error": "Failed to reload Prometheus"})
		return
	}

	log.InfoLog("Pod monitor registered successfully")
	c.JSON(200, gin.H{"message": "pod monitor registered successfully"})

}

func DeletePodMonitor(c *gin.Context) {
	var monitorPod apiObject.MonitorPod
	if err := c.ShouldBindJSON(&monitorPod); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 1. 读取prometheus配置文件
	config := GetPrometheusConfig()
	if config == nil {
		c.JSON(500, gin.H{"error": "Failed to get Prometheus config"})
		return
	}

	// 2. 在制定的job后面把同属于的pod的全部删除，默认“pods”job是用来监控所有pod的
	for j, scrapeConfig := range config.ScrapeConfigs {
		if scrapeConfig.JobName == "pods" {
			for i, uri := range scrapeConfig.StaticConfigs {
				if uri.Labels["instance"] == monitorPod.PodName {
					// 删除这一个instance
					scrapeConfig.StaticConfigs = append(scrapeConfig.StaticConfigs[:i], scrapeConfig.StaticConfigs[i+1:]...)
					break
				}
			}
			if len(scrapeConfig.StaticConfigs) == 0 {
				// 如果这个job没有instance了，就删除这个job
				config.ScrapeConfigs = append(config.ScrapeConfigs[:j], config.ScrapeConfigs[j+1:]...)
			}
		}
	}

	// 3. 保存并写入配置文件
	err := PutPrometheusConfig(config)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to write Prometheus config"})
		return
	}

	// 4. 热加载prometheus
	// curl -X POST http://192.168.1.7:9090/-/reload
	err = exec.Command("curl", "-X", "POST", "http://192.168.1.7:9090/-/reload").Run()
	if err != nil {
		log.ErrorLog("Failed to reload Prometheus: " + err.Error())
		c.JSON(500, gin.H{"error": "Failed to reload Prometheus"})
		return
	}

	log.InfoLog("Pod monitor delete successfully")
	c.JSON(200, gin.H{"message": "pod monitor delete successfully"})
}
