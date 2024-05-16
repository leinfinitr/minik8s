package handlers

import (
	"github.com/gin-gonic/gin"
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

// GetServiceStatus 获取指定Service的状态
func GetServiceStatus(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")

	println("GetServiceStatus: " + namespace + "/" + name)
}
