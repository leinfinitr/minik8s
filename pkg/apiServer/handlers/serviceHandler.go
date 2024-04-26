package handlers

import (
	"github.com/gin-gonic/gin"
)

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
