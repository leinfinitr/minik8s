package handlers

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiObject"
	etcdclient "minik8s/pkg/apiServer/etcdClient"
	"minik8s/pkg/config"
	httprequest "minik8s/tools/httpRequest"
	"minik8s/tools/log"
	"strings"

	"github.com/gin-gonic/gin"
)

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
		// 节点已经存在，需要对pod进行特殊处理，与kubelet同步pod的信息
		log.InfoLog("CreateNode: node already exists")
		c.JSON(config.HttpSuccessCode, "message: node already exists")
		res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdPodPrefix)
		if err != nil {
			log.WarnLog("CreateNode: " + err.Error())
			c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
			return
		}
		var pods []apiObject.Pod
		for _, v := range res {
			var pod apiObject.Pod
			err = json.Unmarshal([]byte(v), &pod)
			if err != nil {
				log.WarnLog("CreateNode: " + err.Error())
				continue
			}
			if pod.Spec.NodeName == node.Metadata.Name {
				pods = append(pods, pod)
			}
		}
		// 把pods信息发送到给kubelet，同步pods信息
		url := "http://" + node.Status.Addresses[0].Address + ":" + fmt.Sprint(config.KubeletAPIPort) + config.PodsSyncURI
		resp, err := httprequest.PostObjMsg(url, pods)
		if err != nil || resp.StatusCode != config.HttpSuccessCode {
			log.WarnLog("CreateNode: " + err.Error())
			c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
			return
		}
		c.JSON(config.HttpSuccessCode, "message: create node success")
		return
	}

	// 注册monitor
	url := config.APIServerURL() + config.MonitorURL
	if resp, err := httprequest.PutObjMsg(url, node); err != nil || resp.StatusCode != config.HttpSuccessCode {
		log.WarnLog("CreateNode: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}

	// 节点首次注册，直接保存节点信息
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
	// 将信息广播给所有node
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
	log.InfoLog("DeleteNodes")
}

// GetNode 获取指定节点
func GetNode(c *gin.Context) {
	name := c.Param("name")
	log.InfoLog("GetNode: " + name)
	// for k, v := range nodes {
	// 	if k == name {
	// 		c.JSON(config.HttpSuccessCode, v)
	// 		return
	// 	}
	// }
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
	// nodes[name] = node

	log.InfoLog("UpdateNode: " + name)
	c.JSON(config.HttpSuccessCode, "")
}

// DeleteNode 删除指定节点
func DeleteNode(c *gin.Context) {
	name := c.Param("name")
	// delete(nodes, name)

	log.InfoLog("DeleteNode: " + name)
	c.JSON(config.HttpSuccessCode, "")
}

// GetNodeStatus 获取指定节点的状态
func GetNodeStatus(c *gin.Context) {
	name := c.Param("name")
	log.InfoLog("GetNodeStatus: " + name)
}

// PingNodeStatus 更新指定节点的状态，其实就是试一试能不能联通
func PingNodeStatus(c *gin.Context) {
	var node apiObject.Node
	err := c.ShouldBindJSON(&node)
	if err != nil {
		log.ErrorLog("PingNodeStatus error: " + err.Error())
		return
	}

	log.DebugLog("start ping NodeIP: " + node.Status.Addresses[0].Address)

	// 尝试三次，失败则认为节点不可用
	times := 0
	success := false
	for times < 3 {
		url := "http://" + node.Status.Addresses[0].Address + ":" + fmt.Sprint(config.KubeletAPIPort) + config.NodeStatusURI
		url = strings.Replace(url, config.NameReplace, node.Metadata.Name, -1)
		resp, err := httprequest.GetMsg(url)
		if err != nil || resp.StatusCode != config.HttpSuccessCode {
			// 无法联通，说明节点不可用
			log.WarnLog("PingNodeStatus: Node can't be connected")
			c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		} else {
			success = true
			newNodeStatus := apiObject.NodeStatus{}
			err = json.NewDecoder(resp.Body).Decode(&newNodeStatus)
			if err != nil {
				log.ErrorLog("PingNodeStatus: " + err.Error())
				c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
				return
			}
			node.Status = newNodeStatus
			nodeJSON, err := json.Marshal(node)
			if err != nil {
				log.ErrorLog("PingNodeStatus: " + err.Error())
				c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
				return
			}
			etcdclient.EtcdStore.Put(config.EtcdNodePrefix+"/"+node.Metadata.Name, string(nodeJSON))
			break
		}
		times++
	}

	if success {
		log.DebugLog("Ping Node success, NodeIp : " + node.Status.Addresses[0].Address)
		c.JSON(config.HttpSuccessCode, "")
	} else {
		log.WarnLog("Ping Node failed : " + node.Status.Addresses[0].Address + " is not available")
		// 无法联通，说明节点不可用
		// 删除monitor配置
		url := config.APIServerURL() + config.MonitorURL
		resp, err := httprequest.DelMsg(url, node)
		if err != nil || resp.StatusCode != config.HttpSuccessCode {
			log.ErrorLog("PingNodeStatus failed")
			c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
			return
		}
		//删除该节点信息
		err = etcdclient.EtcdStore.Delete(config.EtcdNodePrefix + "/" + node.Metadata.Name)
		if err != nil {
			log.ErrorLog("PingNodeStatus failed")
			c.JSON(config.HttpSuccessCode, "")
			return
		}
		log.WarnLog("PingNodeStatus: Node can't be connected, delete node success")
		c.JSON(config.HttpSuccessCode, "")
	}
}
