package handlers

import (
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"
	etcdclient "minik8s/pkg/apiServer/etcdClient"
	"encoding/json"
	"minik8s/pkg/klog"
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
		c.JSON(config.HttpErrorCode,gin.H{"error":err.Error()})
		return
	}
	res,err := etcdclient.EtcdStore.PrefixGet(config.EtcdNodePrefix+"/"+node.Metadata.Name)
	if err != nil {
		klog.WarnLog("APIServer", "CreateNode: "+err.Error())
		c.JSON(config.HttpErrorCode,gin.H{"error":err.Error()})
		return
	}
	if len(res) > 0 {
		klog.WarnLog("APIServer", "CreateNode: node already exists")
		c.JSON(config.HttpErrorCode,gin.H{"error":"node already exists"})
		return
	}
	if node.Kind != apiObject.NodeType {
		klog.WarnLog("APIServer", "CreateNode: node kind is not correct")
		c.JSON(config.HttpErrorCode,gin.H{"error":"node kind is not correct"})
		return
	}
	resJson,err := json.Marshal(node)
	if err != nil {
		klog.WarnLog("APIServer", "CreateNode: "+err.Error())
		c.JSON(config.HttpErrorCode,gin.H{"error":err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(config.EtcdNodePrefix+"/"+node.Metadata.Name,string(resJson))
	if err != nil {
		klog.WarnLog("APIServer", "CreateNode: "+err.Error())
		c.JSON(config.HttpErrorCode,gin.H{"error":err.Error()})
		return
	}
	if len(res)!=1{
		klog.WarnLog("APIServer", "CreateNode: "+err.Error())
		c.JSON(config.HttpErrorCode,gin.H{"error":err.Error()})
		return
	}
	
	// nodes[node.Metadata.Name] = node
	klog.InfoLog("APIServer", "CreateNode: "+node.Metadata.Name)
	c.JSON(config.HttpSuccessCode, "message: create node success")
	//TODO: 将信息广播给所有node
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
	log.DebugLog("UpdateNodeStatus: " + name)

	var status apiObject.NodeStatus
	err := c.ShouldBindJSON(&status)
	if err != nil {
		log.ErrorLog("UpdateNodeStatus error: " + err.Error())

	}

	node := nodes[name]
	node.UpdateNodeStatus(status)
	c.JSON(config.HttpSuccessCode, "")
}
