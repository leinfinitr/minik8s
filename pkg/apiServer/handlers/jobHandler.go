package handlers

import (
	"encoding/json"
	"minik8s/pkg/apiObject"
	etcdclient "minik8s/pkg/apiServer/etcdClient"
	"minik8s/pkg/config"
	"minik8s/tools/log"
	"minik8s/tools/stringops"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetJobs(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	}
	log.InfoLog("GetJobs: " + namespace)
	res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdJobPrefix + "/" + namespace)
	if err != nil {
		log.ErrorLog("GetJobs: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var jobs []apiObject.Job
	for _, v := range res {
		var job apiObject.Job
		err = json.Unmarshal([]byte(v), &job)
		if err != nil {
			log.ErrorLog("GetJobs: " + err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		jobs = append(jobs, job)
	}
	c.JSON(200, jobs)
}

func GetJob(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	} else if name == "" {
		log.ErrorLog("GetJob name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	log.InfoLog("GetJob: " + namespace + "/" + name)
	key := config.EtcdJobPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("GetJob: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var resJson apiObject.Job
	err = json.Unmarshal([]byte(res), &resJson)
	if err != nil {
		log.ErrorLog("GetJob: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, resJson)
	log.InfoLog("GetJob: " + namespace + "/" + name)
}

func AddJob(c *gin.Context) {
	var job apiObject.Job
	err := c.BindJSON(&job)
	if err != nil {
		log.ErrorLog("AddJob: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	newName := job.Metadata.Name
	newNamespace := job.Metadata.Namespace
	if newName == "" {
		log.ErrorLog("AddJob: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	} else if newNamespace == "" {
		job.Metadata.Namespace = "default"
	}
	key := config.EtcdJobPrefix + "/" + job.Metadata.Namespace + "/" + job.Metadata.Name

	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil || len(res) > 0 {
		log.ErrorLog("AddJob: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	job.Metadata.UUID = uuid.New().String()
	randSuffix := "-" + stringops.GenerateRandomString(5)
	reg := regexp.MustCompile(`(.*)(\.[^.]+)`)
	match := reg.FindStringSubmatch(job.Spec.OutputFile)
	if len(match) == 3 {
		job.Spec.OutputFile = match[1] + randSuffix + match[2]
	} else {
		log.ErrorLog("AddJob: output file name is invalid")
		c.JSON(400, gin.H{"error": "output file name is invalid"})
	}
	match = reg.FindStringSubmatch(job.Spec.ErrorFile)
	if len(match) == 3 {
		job.Spec.ErrorFile = match[1] + randSuffix + match[2]
	} else {
		log.ErrorLog("AddJob: error file name is invalid")
		c.JSON(400, gin.H{"error": "error file name is invalid"})
	}
	resJson, err := json.Marshal(job)
	if err != nil {
		log.ErrorLog("AddJob: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(key, string(resJson))
	if err != nil {
		log.ErrorLog("AddJob: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, job)
}

func DeleteJob(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		log.ErrorLog("DeleteJob: name or namespace is empty")
		c.JSON(400, gin.H{"error": "name or namespace is empty"})
		return
	}
	log.InfoLog("DeleteJob: " + namespace + "/" + name)
	key := config.EtcdJobPrefix + "/" + namespace + "/" + name
	err := etcdclient.EtcdStore.Delete(key)
	if err != nil {
		log.ErrorLog("DeleteJob: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": "success"})
	log.InfoLog("DeleteJob: " + namespace + "/" + name + " success")
}

func GetJobStatus(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		log.ErrorLog("GetJobStatus: name or namespace is empty")
		c.JSON(400, gin.H{"error": "name or namespace is empty"})
		return
	}
	log.InfoLog("GetJobStatus: " + namespace + "/" + name)
	key := config.EtcdJobPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("GetJobStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	jb := &apiObject.Job{}
	err = json.Unmarshal([]byte(res), jb)
	if err != nil {
		log.ErrorLog("GetJobStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	status := jb.Status
	if err != nil {
		log.ErrorLog("GetJobStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, status)
}

func UpdateJobStatus(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		log.ErrorLog("UpdateJobStatus: name or namespace is empty")
		c.JSON(400, gin.H{"error": "name or namespace is empty"})
		return
	}
	log.InfoLog("UpdateJobStatus: " + namespace + "/" + name)
	var jobStatus apiObject.JobStatus
	err := c.BindJSON(&jobStatus)
	if err != nil {
		log.ErrorLog("UpdateJobStatus: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	key := config.EtcdJobPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("UpdateJobStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	jb := &apiObject.Job{}
	err = json.Unmarshal([]byte(res), jb)
	if err != nil {
		log.ErrorLog("UpdateJobStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	jb.Status = jobStatus
	resJson, err := json.Marshal(jb)
	if err != nil {
		log.ErrorLog("UpdateJobStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(key, string(resJson))
	if err != nil {
		log.ErrorLog("UpdateJobStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}

func AddJobCode(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		log.ErrorLog("AddJobCode: name or namespace is empty")
		c.JSON(400, gin.H{"error": "name or namespace is empty"})
		return
	}
	log.InfoLog("AddJobCode: " + namespace + "/" + name)
	var code apiObject.JobCode
	err := c.BindJSON(&code)
	if err != nil {
		log.ErrorLog("AddJobCode: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	key := config.EtcdJobCodePrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("AddJobCode: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if len(res) > 0 {
		log.ErrorLog("AddJobCode: " + "code already exists")
		c.JSON(400, gin.H{"error": "code already exists"})
		return
	}
	code.Metadata.UUID = uuid.New().String()
	resJson, err := json.Marshal(code)
	if err != nil {
		log.ErrorLog("AddJobCode: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(key, string(resJson))
	if err != nil {
		log.ErrorLog("AddJobCode: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": string(resJson)})
}

func GetJobCode(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		log.ErrorLog("GetJobCode: name or namespace is empty")
		c.JSON(400, gin.H{"error": "name or namespace is empty"})
		return
	}
	log.InfoLog("GetJobCode: " + namespace + "/" + name)
	key := config.EtcdJobCodePrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("GetJobCode: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var code apiObject.JobCode
	err = json.Unmarshal([]byte(res), &code)
	if err != nil {
		log.ErrorLog("GetJobCode: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, code)
}

func UpdateJobCode(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" || name == "" {
		log.ErrorLog("UpdateJobCode: name or namespace is empty")
		c.JSON(400, gin.H{"error": "name or namespace is empty"})
		return
	}
	log.InfoLog("UpdateJobCode: " + namespace + "/" + name)
	var code apiObject.JobCode
	err := c.BindJSON(&code)
	if err != nil {
		log.ErrorLog("UpdateJobCode: " + err.Error())
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	key := config.EtcdJobCodePrefix + "/" + namespace + "/" + name
	resJson, err := json.Marshal(code)
	if err != nil {
		log.ErrorLog("UpdateJobCode: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(key, string(resJson))
	if err != nil {
		log.ErrorLog("UpdateJobCode: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "success"})
}
