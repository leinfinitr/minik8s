package handlers

import (
	"fmt"
	"minik8s/pkg/apiObject"
	Config "minik8s/pkg/config"
	"minik8s/tools/log"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
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

// 在首次注册时，在本地节点的prometheus配置文件中追加一个监控项
func RegisterMonitor(c *gin.Context) {
	log.InfoLog("Start RegisterMonitor")

	// 1.从请求中获取Node的值
	var node apiObject.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 2. 获取node的IP
	nodeIP := node.Status.Addresses[0].Address
	if nodeIP == Config.APIServerLocalAddress {
		// 如果是apiserver节点，则不需要监控
		c.JSON(200, gin.H{"message": "Node is apiserver"})
		return
	}

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
		log.WarnLog("Node already registered")
		c.JSON(200, gin.H{"message": "Node already registered"})
		return
	}

	// 4. 保存并写入配置文件
	err := PutPrometheusConfig(config)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to write Prometheus config"})
		return
	}

	// 5. 重启prometheus
	// curl -X POST http://192.168.1.7:9090/-/reload
	exec.Command("curl", "-X", "POST", "http://192.168.1.7:9090/-/reload").Run()

	c.JSON(200, gin.H{"message": "monitor registered successfully"})
}

func DeleteMonitor(c *gin.Context) {
	// 1.从请求中获取Node的值
	var node apiObject.Node
	if err := c.ShouldBindJSON(&node); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 2. 获取node的IP
	nodeIP := node.Status.Addresses[0].Address
	if nodeIP == Config.APIServerLocalAddress {
		// 如果是apiserver节点，则不需要监控
		c.JSON(200, gin.H{"message": "Node is apiserver"})
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

	// 5. 重启prometheus
	// curl -X POST http://192.168.1.7:9090/-/reload
	exec.Command("curl", "-X", "POST", "http://192.168.1.7:9090/-/reload").Run()

	c.JSON(200, gin.H{"message": "Node unregistered successfully"})

}
