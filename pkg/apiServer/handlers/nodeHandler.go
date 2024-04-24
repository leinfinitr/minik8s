package handlers

import (
	"github.com/gin-gonic/gin"
)

// GetNodes 获取所有节点
func GetNodes(c *gin.Context) {

	println("GetNodes")
}

// CreateNode 创建节点
func CreateNode(c *gin.Context) {

	println("CreateNode")
}

// DeleteNodes 删除所有节点
func DeleteNodes(c *gin.Context) {

	println("DeleteNodes")
}

// GetNode 获取指定节点
func GetNode(c *gin.Context) {
	name := c.Param("name")

	println("GetNode: " + name)
}

// UpdateNode 更新指定节点
func UpdateNode(c *gin.Context) {
	name := c.Param("name")

	println("UpdateNode: " + name)

}

// DeleteNode 删除指定节点
func DeleteNode(c *gin.Context) {
	name := c.Param("name")

	println("DeleteNode: " + name)
}

// GetNodeStatus 获取指定节点的状态
func GetNodeStatus(c *gin.Context) {
	name := c.Param("name")

	println("GetNodeStatus: " + name)
}

// UpdateNodeStatus 更新指定节点的状态
func UpdateNodeStatus(c *gin.Context) {
	name := c.Param("name")

	println("UpdateNodeStatus: " + name)
}
