package handlers

import (
	"encoding/json"
	"minik8s/pkg/apiObject"
	etcdclient "minik8s/pkg/apiServer/etcdClient"
	"minik8s/pkg/config"
	"minik8s/tools/log"

	"github.com/gin-gonic/gin"
)


func GetDNS(c *gin.Context){
	name := c.Param("name")
	namespace := c.Param("namespace")
	if name == ""{
		log.ErrorLog("GetDNS name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	if namespace == ""{
		namespace = "default"
	}
	log.InfoLog("GetDNS: "+namespace+"/"+name)
	
	key := config.EtcdDnsPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil{
		log.ErrorLog("GetDNS: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(res) == 0{
		log.ErrorLog("GetDNS: not found")
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	var resJson apiObject.Dns
	err = json.Unmarshal([]byte(res), &resJson)
	if err != nil{
		log.ErrorLog("GetDNS: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": resJson})
}

func GetDNSs(c *gin.Context){
	namespace := c.Param("namespace")
	if namespace == ""{
		namespace = "default"
	}
	log.InfoLog("GetDNSs: "+namespace)
	key := config.EtcdDnsPrefix + "/" + namespace
	res, err := etcdclient.EtcdStore.PrefixGet(key)
	if err != nil{
		log.ErrorLog("GetDNSs: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var dnss []apiObject.Dns
	for _,v := range res{
		var dns apiObject.Dns
		err = json.Unmarshal([]byte(v), &dns)
		if err != nil{
			log.ErrorLog("GetDNSs: "+err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		dnss = append(dnss, dns)
	}
	c.JSON(200, gin.H{"data": dnss})
}

func AddDNS(c *gin.Context){
	var dns apiObject.Dns
	err := c.BindJSON(&dns)
	if err != nil{
		log.ErrorLog("AddDNS: "+err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	log.InfoLog("AddDNS: "+dns.Metadata.Namespace+"/"+dns.Metadata.Name)
	
	key := config.EtcdDnsPrefix + "/" + dns.Metadata.Namespace + "/" + dns.Metadata.Name
	oldRes, err := etcdclient.EtcdStore.Get(key)
	if err != nil{
		log.ErrorLog("AddDNS: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(oldRes) != 0{
		log.ErrorLog("AddDNS: already exists")
		c.JSON(409, gin.H{"error": "already exists"})
		return
	}
	if dns.Metadata.Name == ""{
		log.ErrorLog("AddDNS: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	if dns.Metadata.Namespace == ""{
		dns.Metadata.Namespace = "default"
	}

	//TODO
}

func DeleteDNS(c *gin.Context){
	name := c.Param("name")
	namespace := c.Param("namespace")
	if name == ""{
		log.ErrorLog("DeleteDNS name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	if namespace == ""{
		namespace = "default"
	}
	log.InfoLog("DeleteDNS: "+namespace+"/"+name)
	
	key := config.EtcdDnsPrefix + "/" + namespace + "/" + name
	oldRes, err := etcdclient.EtcdStore.Get(key)
	if err != nil{
		log.ErrorLog("DeleteDNS: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(oldRes) == 0{
		log.ErrorLog("DeleteDNS: not found")
		c.JSON(404, gin.H{"error": "not found"})
		return
	}
	err = etcdclient.EtcdStore.Delete(key)
	if err != nil{
		log.ErrorLog("DeleteDNS: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	//TODO dnsUpdate sync
	c.JSON(200, gin.H{"data": "success"})
}