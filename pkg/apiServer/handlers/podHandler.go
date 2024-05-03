package handlers

import (
	etcdclient "minik8s/pkg/apiServer/etcdClient"
	"minik8s/pkg/config"
	klog "minik8s/pkg/klog"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/apiObject"
)

// GetPod 获取指定Pod
func GetPod(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	}else if name == "" {
		klog.WarnLog("APIServer", "GetPod: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	klog.InfoLog("APIServer", "GetPod: "+namespace+"/"+name)
	key := config.EtcdPodPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		klog.WarnLog("APIServer", "GetPod: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resJson,err := json.Marshal(res)
	if err != nil {
		klog.WarnLog("APIServer", "GetPod: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": resJson})

}

// UpdatePod 更新Pod
func UpdatePod(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == ""{
		namespace = "default"
	}else if name == "" {
		klog.WarnLog("APIServer", "UpdatePod: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	klog.InfoLog("APIServer", "UpdatePod: "+namespace+"/"+name)
	key := config.EtcdPodPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		klog.WarnLog("APIServer", "UpdatePod: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resJson,err := json.Marshal(res)
	if err != nil {
		klog.WarnLog("APIServer", "UpdatePod: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	pod := &apiObject.Pod{}
	//解析json到pod
	err = json.Unmarshal(resJson,pod)
	if err != nil {
		klog.WarnLog("APIServer", "UpdatePod: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	reqPod := &apiObject.Pod{}
	err = c.ShouldBindJSON(reqPod)
	if err != nil {
		klog.WarnLog("APIServer", "UpdatePod: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	//更新pod
	UpdatePodProps(pod,reqPod)
	//将更新后的pod写入etcd
	resJson,err = json.Marshal(pod)
	if err != nil {
		klog.WarnLog("APIServer", "UpdatePod: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if pod.Metadata.Namespace == "" || pod.Metadata.Name == "" {
		klog.WarnLog("APIServer", "UpdatePod: namespace or name is empty")
		c.JSON(400, gin.H{"error": "namespace or name is empty"})
		return
	}
	key = config.EtcdPodPrefix + "/" + pod.Metadata.Namespace + "/" + pod.Metadata.Name
	err = etcdclient.EtcdStore.Put(key,string(resJson))
	if err != nil {
		klog.WarnLog("APIServer", "UpdatePod: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": resJson})
}

// DeletePod 删除Pod
func DeletePod(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")

	println("DeletePod: " + namespace + "/" + name)
}

// GetPodEphemeralContainers 获取指定Pod的临时容器
func GetPodEphemeralContainers(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")

	println("GetPodEphemeralContainers: " + namespace + "/" + name)
}

// UpdatePodEphemeralContainers 更新指定Pod的临时容器
func UpdatePodEphemeralContainers(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")

	println("UpdatePodEphemeralContainers: " + namespace + "/" + name)
}

// GetPodLog 获取指定Pod的日志
func GetPodLog(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")

	println("GetPodLog: " + namespace + "/" + name)
}

// GetPodStatus 获取指定Pod的状态
func GetPodStatus(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")

	println("GetPodStatus: " + namespace + "/" + name)
}

// UpdatePodStatus 更新指定Pod的状态
func UpdatePodStatus(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")

	println("UpdatePodStatus: " + namespace + "/" + name)
}

// GetPods 获取所有Pod
func GetPods(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	}
	klog.InfoLog("APIServer", "GetPods: "+namespace)
	key := config.EtcdPodPrefix + "/" + namespace
	res, err := etcdclient.EtcdStore.PrefixGet(key)
	if err != nil {
		klog.WarnLog("APIServer", "GetPods: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resJson,err := json.Marshal(res)
	if err != nil {
		klog.WarnLog("APIServer", "GetPods: "+err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": resJson})
}

// CreatePod 创建Pod
func CreatePod(c *gin.Context) {
	namespace := c.Param("namespace")

	println("CreatePod: " + namespace)
}

// DeletePods 删除所有Pod
func DeletePods(c *gin.Context) {
	namespace := c.Param("namespace")

	println("DeletePods: " + namespace)
}

// GetGlobalPods 获取全局所有Pod
func GetGlobalPods(c *gin.Context) {
	println("GetGlobalPods")
}


//TODO: 更新Pod
func UpdatePodProps(old *apiObject.Pod, new *apiObject.Pod) {
	
}
