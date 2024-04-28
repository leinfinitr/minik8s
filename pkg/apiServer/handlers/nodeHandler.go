package handlers

import (
	"github.com/gin-gonic/gin"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"
)

// 一个临时的用于存储 node 信息的 map
var nodes = make(map[string]apiObject.Node)

// GetNodes 获取所有节点
func GetNodes(c *gin.Context) {

	log.InfoLog("GetNodes")
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
	c.JSON(config.SuccessCode, "")
}

// DeleteNodes 删除所有节点
func DeleteNodes(c *gin.Context) {

	log.InfoLog("DeleteNodes")
}

// GetNode 获取指定节点
func GetNode(c *gin.Context) {
	name := c.Param("name")

	log.InfoLog("GetNode: " + name)
	c.JSON(200, nodes[name])
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
	c.JSON(200, nodes[name])
}

// DeleteNode 删除指定节点
func DeleteNode(c *gin.Context) {
	name := c.Param("name")

	log.InfoLog("DeleteNode: " + name)
}

// GetNodeStatus 获取指定节点的状态
func GetNodeStatus(c *gin.Context) {
	name := c.Param("name")

	log.InfoLog("GetNodeStatus: " + name)
}

// UpdateNodeStatus 更新指定节点的状态
func UpdateNodeStatus(c *gin.Context) {
	name := c.Param("name")

	log.InfoLog("UpdateNodeStatus: " + name)
}
