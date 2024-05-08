package handlers

import (
	"encoding/json"
	"fmt"
	httprequest "minik8s/internal/pkg/httpRequest"
	"minik8s/pkg/apiObject"
	etcdclient "minik8s/pkg/apiServer/etcdClient"
	"minik8s/pkg/config"
	"minik8s/tools/log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetPod 获取指定Pod
func GetPod(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	} else if name == "" {
		log.WarnLog("GetPod: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	log.InfoLog("GetPod: " + namespace + "/" + name)
	key := config.EtcdPodPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.WarnLog("GetPod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resJson, err := json.Marshal(res)
	if err != nil {
		log.WarnLog("GetPod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": resJson})

}

// UpdatePod 更新Pod
func UpdatePod(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	} else if name == "" {
		log.WarnLog("UpdatePod: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	log.InfoLog("UpdatePod: " + namespace + "/" + name)
	key := config.EtcdPodPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.WarnLog("UpdatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resJson, err := json.Marshal(res)
	if err != nil {
		log.WarnLog("UpdatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	pod := &apiObject.Pod{}
	//解析json到pod
	err = json.Unmarshal(resJson, pod)
	if err != nil {
		log.WarnLog("UpdatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	reqPod := &apiObject.Pod{}
	err = c.ShouldBindJSON(reqPod)
	if err != nil {
		log.WarnLog("UpdatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	//更新pod
	UpdatePodProps(pod, reqPod)
	//将更新后的pod写入etcd
	resJson, err = json.Marshal(pod)
	if err != nil {
		log.WarnLog("UpdatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if pod.Metadata.Namespace == "" || pod.Metadata.Name == "" {
		log.WarnLog("UpdatePod: namespace or name is empty")
		c.JSON(400, gin.H{"error": "namespace or name is empty"})
		return
	}
	key = config.EtcdPodPrefix + "/" + pod.Metadata.Namespace + "/" + pod.Metadata.Name
	err = etcdclient.EtcdStore.Put(key, string(resJson))
	if err != nil {
		log.WarnLog("UpdatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": resJson})
}

// DeletePod 删除Pod
func DeletePod(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" {
		namespace = "default"
	} else if name == "" {
		log.WarnLog("DeletePods: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	key := config.EtcdPodPrefix + "/" + namespace + "/" + name
	log.WarnLog("DeletePods: " + namespace + "/" + name)
	err := etcdclient.EtcdStore.Delete(key)
	if err != nil {
		log.WarnLog("DeletePods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	//TODO删除kube-controller-manager中的Endpoint
	//TODO发送删除请求到kubelet
	url := config.KubeletLocalIP + ":" + fmt.Sprint(config.KubeletAPIPort)
	delUri := url + config.PodsURI
	delUri = strings.Replace(delUri, config.NameSpaceReplace, namespace, -1)
	delUri = strings.Replace(delUri, config.NameReplace, name, -1)
	_, err = httprequest.DelObjMsg(delUri)
	if err != nil {
		log.WarnLog("DeletePods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": "delete success"})
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
	if namespace == "" {
		namespace = "default"
	}else if name == "" {
		log.WarnLog("GetPodStatus: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	key := config.EtcdPodPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.WarnLog("GetPodStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resJson, err := json.Marshal(res)
	if err != nil {
		log.WarnLog("GetPodStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	pod	:= &apiObject.Pod{}
	err = json.Unmarshal(resJson, pod)
	if err != nil {
		log.WarnLog("GetPodStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	//获取pod状态
	status := pod.Status
	byteStatus, err := json.Marshal(status)
	if err != nil {
		log.WarnLog("GetPodStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": byteStatus})

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
	log.InfoLog("GetPods: " + namespace)
	key := config.EtcdPodPrefix + "/" + namespace
	res, err := etcdclient.EtcdStore.PrefixGet(key)
	if err != nil {
		log.WarnLog("GetPods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resJson, err := json.Marshal(res)
	if err != nil {
		log.WarnLog("GetPods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": resJson})
}

// CreatePod 创建Pod
func CreatePod(c *gin.Context) {
	pod := &apiObject.Pod{}
	err := c.ShouldBindJSON(pod)
	if err != nil {
		log.WarnLog("CreatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	newPodName := pod.Metadata.Name
	newPodNamespace := pod.Metadata.Namespace
	if newPodName == "" || newPodNamespace == "" {
		log.WarnLog("CreatePod: name or namespace is empty")
		c.JSON(400, gin.H{"error": "name or namespace is empty"})
		return
	}
	key := config.EtcdPodPrefix + "/" + pod.Metadata.Namespace + "/" + pod.Metadata.Name
	_, err = etcdclient.EtcdStore.Get(key)
	if err == nil {
		log.WarnLog("CreatePod: Pod already exists")
		c.JSON(400, gin.H{"error": "Pod already exists"})
		return
	}
	log.InfoLog("CreatePod: " + newPodNamespace + "/" + newPodName)
	// TODO: 生成 UUID
	pod.Metadata.UUID = uuid.New().String()
	// TODO: 发送的时候筛选 node
	url := config.KubeletLocalIP + ":" + fmt.Sprint(config.KubeletAPIPort)
	createUri := url + config.PodsURI
	createUri = strings.Replace(createUri, config.NameSpaceReplace, newPodNamespace, -1)
	createUri = strings.Replace(createUri, config.NameReplace, newPodName, -1)
	reaJson, err := json.Marshal(pod)
	if err != nil {
		log.WarnLog("CreatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(key, string(reaJson))
	if err != nil {
		log.WarnLog("CreatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resp, err := httprequest.PostObjMsg(createUri, pod)
	if err != nil {
		fmt.Println("Error: Could not post the object message.")
		os.Exit(1)
	}
	c.JSON(200, gin.H{"data": resp})
}

// DeletePods 删除所有Pod
func DeletePods(c *gin.Context) {
	namespace := c.Param("namespace")

	println("DeletePods: " + namespace)
}

// GetGlobalPods 获取全局所有Pod
func GetGlobalPods(c *gin.Context) {
	key := config.EtcdPodPrefix
	res, err := etcdclient.EtcdStore.PrefixGet(key)
	if err != nil {
		log.WarnLog("GetGlobalPods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resJson, err := json.Marshal(res)
	if err != nil {
		log.WarnLog("GetGlobalPods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": resJson})
}

// TODO: 更新Pod
func UpdatePodProps(old *apiObject.Pod, new *apiObject.Pod) {

}
