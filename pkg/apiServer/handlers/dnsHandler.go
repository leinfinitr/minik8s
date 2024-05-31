package handlers

import (
	"encoding/json"
	"minik8s/pkg/apiObject"
	etcdclient "minik8s/pkg/apiServer/etcdClient"
	"minik8s/pkg/config"
	"minik8s/tools/log"
	"github.com/google/uuid"
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

	for it,path := range dns.Spec.Paths{
		log.InfoLog("AddDNS: "+dns.Metadata.Namespace+"/"+dns.Metadata.Name+" path "+string(it))
		if path.SvcName == ""{
			log.ErrorLog("AddDNS: svcName is empty")
			c.JSON(400, gin.H{"error": "svcName is empty"})
			return
		}
		svcKey := config.EtcdServicePrefix + "/" + dns.Metadata.Namespace + "/" + path.SvcName
		svcRes, err := etcdclient.EtcdStore.Get(svcKey)
		if err != nil{
			log.ErrorLog("AddDNS: "+err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if len(svcRes) == 0{
			log.ErrorLog("AddDNS: svc not found")
			c.JSON(404, gin.H{"error": "svc not found"})
			return
		}
		service := apiObject.Service{}
		err = json.Unmarshal([]byte(svcRes), &service)
		if err != nil{
			log.ErrorLog("AddDNS: "+err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		dns.Spec.Paths[it].SvcIp = service.Spec.ClusterIP
	}
	dns.Metadata.UUID = uuid.New().String()
	//Update dnsRequest
	var dnsRequest apiObject.DnsRequest
	dnsRequest.Action = "Create"
	dnsRequest.DnsMeta = dns.Metadata
	dnsRequestJson,err := json.Marshal(dnsRequest)
	if err != nil{
		log.ErrorLog("AddDNS: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	key = config.EtcdDnsRequestPrefix + "/" + dns.Metadata.Namespace + "/" + dns.Metadata.Name
	err = etcdclient.EtcdStore.Put(key, string(dnsRequestJson))
	if err != nil{
		log.ErrorLog("AddDNS: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"data": dns})
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
	//Update dnsRequest
	var dns apiObject.Dns
	err = json.Unmarshal([]byte(oldRes), &dns)
	if err != nil{
		log.ErrorLog("DeleteDNS: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Delete(key)
	if err != nil{
		log.ErrorLog("DeleteDNS: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var dnsRequest apiObject.DnsRequest
	dnsRequest.Action = "Delete"
	dnsRequest.DnsMeta = dns.Metadata
	dnsRequestJson,err := json.Marshal(dnsRequest)
	if err != nil{
		log.ErrorLog("DeleteDNS: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	key = config.EtcdDnsRequestPrefix + "/" + namespace + "/" + name
	err = etcdclient.EtcdStore.Put(key, string(dnsRequestJson))
	if err != nil{
		log.ErrorLog("DeleteDNS: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
}