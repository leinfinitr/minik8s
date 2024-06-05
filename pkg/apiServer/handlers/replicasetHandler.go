package handlers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"

	etcdclient "minik8s/pkg/apiServer/etcdClient"
)

func GetReplicaSet(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	} else if name == "" {
		log.ErrorLog("GetReplicaSet name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	log.InfoLog("GetReplicaSet: " + namespace + "/" + name)

	key := config.EtcdReplicaSetPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("GetReplicaSet: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var resJson apiObject.ReplicaSet
	err = json.Unmarshal([]byte(res), &resJson)
	if err != nil {
		log.ErrorLog("GetReplicaSet: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": resJson})
	log.InfoLog("GetReplicaSet: " + namespace + "/" + name)
}

func GetReplicaSets(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	}
	log.InfoLog("GetReplicaSets: " + namespace)

	key := config.EtcdReplicaSetPrefix + "/" + namespace
	res, err := etcdclient.EtcdStore.PrefixGet(key)
	if err != nil {
		log.ErrorLog("GetReplicaSets: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var replicaSets []apiObject.ReplicaSet
	for _, v := range res {
		var rs apiObject.ReplicaSet
		err = json.Unmarshal([]byte(v), &rs)
		if err != nil {
			log.ErrorLog("GetReplicaSets: " + err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		replicaSets = append(replicaSets, rs)
	}
	c.JSON(200, replicaSets)
	log.InfoLog("GetReplicaSets: " + namespace)
}

func AddReplicaSet(c *gin.Context) {
	var rs apiObject.ReplicaSet
	err := c.ShouldBindJSON(&rs)
	if err != nil {
		log.ErrorLog("AddReplicaSet error: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdReplicaSetPrefix + "/" + rs.Metadata.Namespace + "/" + rs.Metadata.Name)
	if err != nil {
		log.ErrorLog("AddReplicaSet: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(res) > 0 {
		log.ErrorLog("AddReplicaSet: replicaSet already exists")
		c.JSON(500, gin.H{"error": "replicaSet already exists"})
		return
	}

	rs.Metadata.UUID = uuid.New().String()
	resJson, err := json.Marshal(rs)
	if err != nil {
		log.ErrorLog("AddReplicaSet: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(config.EtcdReplicaSetPrefix+"/"+rs.Metadata.Namespace+"/"+rs.Metadata.Name, string(resJson))
	if err != nil {
		log.ErrorLog("AddReplicaSet: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": string(resJson)})
	log.InfoLog("AddReplicaSet: " + rs.Metadata.Namespace + "/" + rs.Metadata.Name)
}

func DeleteReplicaSet(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		log.ErrorLog("DeleteReplicaSet: namespace or name is empty")
		c.JSON(400, gin.H{"error": "namespace or name is empty"})
		return
	}
	log.InfoLog("DeleteReplicaSet: " + namespace + "/" + name)

	key := config.EtcdReplicaSetPrefix + "/" + namespace + "/" + name
	err := etcdclient.EtcdStore.Delete(key)
	if err != nil {
		log.ErrorLog("DeleteReplicaSet: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": "success"})
	log.InfoLog("DeleteReplicaSet: " + namespace + "/" + name + " success")
}

func UpdateReplicaSet(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	} else if name == "" {
		log.ErrorLog("UpdateReplicaSet: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	log.InfoLog("UpdateReplicaSet: " + namespace + "/" + name)

	key := config.EtcdReplicaSetPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if res == "" || err != nil {
		log.ErrorLog("UpdateReplicaSet: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	rs := &apiObject.ReplicaSet{}
	err = c.ShouldBindJSON(rs)
	if err != nil {
		log.ErrorLog("UpdateReplicaSet: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resJson, err := json.Marshal(rs)
	if err != nil {
		log.ErrorLog("UpdateReplicaSet: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if rs.Metadata.Namespace == "" || rs.Metadata.Name == "" {
		log.ErrorLog("UpdateReplicaSet: namespace or name is empty")
		c.JSON(400, gin.H{"error": "namespace or name is empty"})
		return
	}
	key = config.EtcdReplicaSetPrefix + "/" + rs.Metadata.Namespace + "/" + rs.Metadata.Name
	err = etcdclient.EtcdStore.Put(key, string(resJson))
	if err != nil {
		log.ErrorLog("UpdateReplicaSet: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": string(resJson)})
	log.InfoLog("UpdateReplicaSet: " + namespace + "/" + name + " success")
}

func GetGlobalReplicaSets(c *gin.Context) {
	log.DebugLog("GetGlobalReplicaSets")
	key := config.EtcdReplicaSetPrefix
	res, err := etcdclient.EtcdStore.PrefixGet(key)
	if err != nil {
		log.ErrorLog("GetGlobalReplicaSets: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var replicaSets []apiObject.ReplicaSet
	for _, v := range res {
		var rs apiObject.ReplicaSet
		err = json.Unmarshal([]byte(v), &rs)
		if err != nil {
			log.ErrorLog("GetGlobalReplicaSets: " + err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		replicaSets = append(replicaSets, rs)
	}
	c.JSON(200, replicaSets)
	log.DebugLog("GetGlobalReplicaSets success")
}

func GetReplicaSetStatus(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	} else if name == "" {
		log.ErrorLog("GetReplicaSetStatus: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	log.InfoLog("GetReplicaSetStatus: " + namespace + "/" + name)
	key := config.EtcdReplicaSetPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("GetReplicaSetStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	rs := &apiObject.ReplicaSet{}
	err = json.Unmarshal([]byte(res), rs)
	if err != nil {
		log.ErrorLog("GetReplicaSetStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	status := rs.Status
	byteStatus, err := json.Marshal(status)
	if err != nil {
		log.ErrorLog("GetReplicaSetStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": byteStatus})
}

func UpdateReplicaSetStatus(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" || name == "" {
		log.ErrorLog("UpdateReplicaSetStatus: namespace or name is empty")
		c.JSON(400, gin.H{"error": "namespace or name is empty"})
		return
	}
	log.DebugLog("UpdateReplicaSetStatus: " + namespace + "/" + name)
	var status apiObject.ReplicaSetStatus
	err := c.ShouldBindJSON(&status)
	if err != nil {
		log.ErrorLog("UpdateReplicaSetStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	key := config.EtcdReplicaSetPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("UpdateReplicaSetStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	rs := &apiObject.ReplicaSet{}
	err = json.Unmarshal([]byte(res), rs)
	if err != nil {
		log.ErrorLog("UpdateReplicaSetStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	rs.Status = status
	resJson, err := json.Marshal(rs)
	if err != nil {
		log.ErrorLog("UpdateReplicaSetStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(key, string(resJson))
	if err != nil {
		log.ErrorLog("UpdateReplicaSetStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": string(resJson)})
	log.DebugLog("UpdateReplicaSetStatus: " + namespace + "/" + name + " success")
}
