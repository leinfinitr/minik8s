package handlers

import (
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"

	"github.com/gin-gonic/gin"
)

// 一个临时的用于存储 node 信息的 map
var nodes = make(map[string]apiObject.Node)

// GetNodes 获取所有节点
func GetNodes(c *gin.Context) {
	log.InfoLog("GetNodes")

	c.JSON(config.HttpSuccessCode, nodes)
}

// CreateNode 创建节点
func CreateNode(c *gin.Context) {
	var node apiObject.Node
	err := c.ShouldBindJSON(&node)
	if err != nil {
		log.ErrorLog("CreateNode error: " + err.Error())
	}
	nodes[node.Metadata.Name] = node

	log.InfoLog("CreateNode: " + node.Metadata.Name)
	c.JSON(config.HttpSuccessCode, "")
}

// DeleteNodes 删除所有节点
func DeleteNodes(c *gin.Context) {
	nodes = make(map[string]apiObject.Node)
	log.InfoLog("DeleteNodes")
}

// GetNode 获取指定节点
func GetNode(c *gin.Context) {
	name := c.Param("name")
	log.InfoLog("GetNode: " + name)
	for k, v := range nodes {
		if k == name {
			c.JSON(config.HttpSuccessCode, v)
			return
		}
	}
	c.JSON(config.HttpNotFoundCode, "")
}

// UpdateNode 更新指定节点
func UpdateNode(c *gin.Context) {
	var node apiObject.Node
	err := c.ShouldBindJSON(&node)
	if err != nil {
		log.ErrorLog("UpdateNode error: " + err.Error())
	}
	name := c.Param("name")
	nodes[name] = node

	log.InfoLog("UpdateNode: " + name)
	c.JSON(config.HttpSuccessCode, "")
}

// DeleteNode 删除指定节点
func DeleteNode(c *gin.Context) {
	name := c.Param("name")
	delete(nodes, name)

	log.InfoLog("DeleteNode: " + name)
	c.JSON(config.HttpSuccessCode, "")
}

// GetNodeStatus 获取指定节点的状态
func GetNodeStatus(c *gin.Context) {
	name := c.Param("name")
	log.InfoLog("GetNodeStatus: " + name)

	for k, v := range nodes {
		if k == name {
			c.JSON(config.HttpSuccessCode, v.Status)
			return
		}
	}

}

// UpdateNodeStatus 更新指定节点的状态
func UpdateNodeStatus(c *gin.Context) {
	name := c.Param("name")
	log.InfoLog("UpdateNodeStatus: " + name)

	var status apiObject.NodeStatus
	err := c.ShouldBindJSON(&status)
	if err != nil {
		log.ErrorLog("UpdateNodeStatus error: " + err.Error())

	}

	node := nodes[name]
	node.UpdateNodeStatus(status)
	c.JSON(config.HttpSuccessCode, "")
}
