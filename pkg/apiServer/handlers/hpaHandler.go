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

func GetHPA(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	}
	log.InfoLog("GetHPA: " + namespace)
	name := c.Param("name")
	if name == "" {
		log.ErrorLog("GetHPA name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	log.InfoLog("GetHPA: " + namespace + "/" + name)

	res, err := etcdclient.EtcdStore.Get(config.EtcdHpaPrefix + "/" + namespace + "/" + name)
	if err != nil {
		log.ErrorLog("GetHPA: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(res) == 0 {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	if len(res) > 1 {
		log.ErrorLog("GetHPA: more than one HPA")
		c.JSON(500, gin.H{"error": "more than one HPA"})
		return
	}
	c.JSON(200, gin.H{"data": res[0]})
}

func GetHPAs(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	}
	log.InfoLog("GetHPAs: " + namespace)
	res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdHpaPrefix + "/" + namespace)
	if err != nil {
		log.ErrorLog("GetHPAs: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var hpaList []apiObject.HPA
	for _, v := range res {
		hpa := apiObject.HPA{}
		err = json.Unmarshal([]byte(v), &hpa)
		if err != nil {
			log.ErrorLog("GetHPAs: " + err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		hpaList = append(hpaList, hpa)
	}

	c.JSON(200, hpaList)
}

func AddHPA(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	}
	log.InfoLog("AddHPA: " + namespace)

	var hpa apiObject.HPA
	err := c.BindJSON(&hpa)
	if err != nil {
		log.ErrorLog("AddHPA: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if hpa.Metadata.Name == "" {
		log.ErrorLog("AddHPA: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	if hpa.Metadata.Namespace == "" {
		hpa.Metadata.Namespace = "default"
	}
	//检查是否已经存在
	key := config.EtcdHpaPrefix + "/" + hpa.Metadata.Namespace + "/" + hpa.Metadata.Name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("AddHPA: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(res) > 0 {
		log.ErrorLog("AddHPA: already exists")
		c.JSON(400, gin.H{"error": "already exists"})
		return
	}
	hpa.Metadata.UUID = uuid.New().String()
	resJson, err := json.Marshal(hpa)
	if err != nil {
		log.ErrorLog("AddHPA: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(key, string(resJson))
	if err != nil {
		log.ErrorLog("AddHPA: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": hpa})
}

func DeleteHPA(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	}
	log.InfoLog("DeleteHPA: " + namespace)

	name := c.Param("name")
	if name == "" {
		log.ErrorLog("DeleteHPA name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	log.InfoLog("DeleteHPA: " + namespace + "/" + name)

	err := etcdclient.EtcdStore.Delete(config.EtcdHpaPrefix + "/" + namespace + "/" + name)
	if err != nil {
		log.ErrorLog("DeleteHPA: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": "success"})
}

func GetGlobalHPAs(c *gin.Context) {
	log.DebugLog("GetGlobalHPAs")
	res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdHpaPrefix)
	if err != nil {
		log.ErrorLog("GetGlobalHPAs: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var hpaList []apiObject.HPA
	for _, v := range res {
		hpa := apiObject.HPA{}
		err = json.Unmarshal([]byte(v), &hpa)
		if err != nil {
			log.ErrorLog("GetGlobalHPAs: " + err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		hpaList = append(hpaList, hpa)
	}
	c.JSON(200, hpaList)
}

func UpdateHPAStatus(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		log.ErrorLog("UpdateHPAStatus: name or namespace is empty")
		c.JSON(400, gin.H{"error": "name or namespace is empty"})
		return
	}
	log.DebugLog("UpdateHPAStatus: " + namespace + "/" + name)

	var status apiObject.HPAStatus
	err := c.BindJSON(&status)
	if err != nil {
		log.ErrorLog("UpdateHPAStatus: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	key := config.EtcdHpaPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("UpdateHPAStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	hpa := apiObject.HPA{}
	err = json.Unmarshal([]byte(res), &hpa)
	if err != nil {
		log.ErrorLog("UpdateHPAStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	hpa.Status = status
	resJson, err := json.Marshal(hpa)
	if err != nil {
		log.ErrorLog("UpdateHPAStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(key, string(resJson))
	if err != nil {
		log.ErrorLog("UpdateHPAStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": string(resJson)})
}
