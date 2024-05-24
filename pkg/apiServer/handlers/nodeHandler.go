package handlers

import (
	"encoding/json"
	"minik8s/pkg/apiObject"
	etcdclient "minik8s/pkg/apiServer/etcdClient"
	"minik8s/pkg/config"
	"minik8s/tools/log"

	"github.com/gin-gonic/gin"
)

// 一个临时的用于存储 node 信息的 map
var nodes = make(map[string]apiObject.Node)

// GetNodes 获取所有节点
func GetNodes(c *gin.Context) {
	log.InfoLog("GetNodes")
	res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdNodePrefix)
	if err != nil {
		log.WarnLog("GetNodes: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}

	var nodes []apiObject.Node
	for _, v := range res {
		var node apiObject.Node
		err = json.Unmarshal([]byte(v), &node)
		if err != nil {
			log.WarnLog("GetNodes: " + err.Error())
			c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
			return
		}
		nodes = append(nodes, node)
	}

	c.JSON(config.HttpSuccessCode, nodes)
}

// CreateNode 创建节点
func CreateNode(c *gin.Context) {
	var node apiObject.Node
	err := c.ShouldBindJSON(&node)
	if err != nil {
		log.ErrorLog("CreateNode error: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}
	res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdNodePrefix + "/" + node.Metadata.Name)
	if err != nil {
		log.WarnLog("CreateNode: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}
	if len(res) > 0 {
		log.WarnLog("CreateNode: node already exists")
		c.JSON(config.HttpErrorCode, gin.H{"error": "node already exists"})
		return
	}
	if node.Kind != apiObject.NodeType {
		log.WarnLog("CreateNode: node kind is not correct")
		c.JSON(config.HttpErrorCode, gin.H{"error": "node kind is not correct"})
		return
	}
	resJson, err := json.Marshal(node)
	if err != nil {
		log.WarnLog("CreateNode: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(config.EtcdNodePrefix+"/"+node.Metadata.Name, string(resJson))
	if err != nil {
		log.WarnLog("CreateNode: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}

	// nodes[node.Metadata.Name] = node
	log.InfoLog("CreateNode: " + node.Metadata.Name + " Node IP: " + node.Status.Addresses[0].Address)
	c.JSON(config.HttpSuccessCode, "message: create node success")
	// TODO: 将信息广播给所有node
	BroadcastNode(node)
}

func BroadcastNode(node apiObject.Node) {
	res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdPodPrefix)
	if err != nil {
		log.WarnLog("BroadcastNode: " + err.Error())
		return
	}
	for _, v := range res {
		var pod apiObject.Pod
		err = json.Unmarshal([]byte(v), &pod)
		if err != nil {
			log.WarnLog("BroadcastNode: " + err.Error())
			continue
		}
		if pod.Spec.NodeName == node.Metadata.Name {
			pod.Status.Phase = apiObject.PodRunning
			resJson, err := json.Marshal(pod)
			if err != nil {
				log.WarnLog("BroadcastNode: " + err.Error())
				continue
			}
			err = etcdclient.EtcdStore.Put(config.EtcdPodPrefix+"/"+pod.Metadata.Name, string(resJson))
			if err != nil {
				log.WarnLog("BroadcastNode: " + err.Error())
				continue
			}
		}
	}
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
