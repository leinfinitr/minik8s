package handlers

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiObject"
	etcdclient "minik8s/pkg/apiServer/etcdClient"
	"minik8s/pkg/config"
	httprequest "minik8s/tools/httpRequest"
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
		log.ErrorLog("GetPod: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	log.InfoLog("GetPod: " + namespace + "/" + name)
	key := config.EtcdPodPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("GetPod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resJson, err := json.Marshal(res)
	if err != nil {
		log.ErrorLog("GetPod: " + err.Error())
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
		log.ErrorLog("UpdatePod: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	log.InfoLog("UpdatePod: " + namespace + "/" + name)
	key := config.EtcdPodPrefix + "/" + namespace + "/" + name
	_, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("UpdatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	pod := &apiObject.Pod{}
	err = c.ShouldBindJSON(pod)
	if err != nil {
		log.ErrorLog("UpdatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 更新pod
	UpdatePodProps(pod)
	// 将更新后的pod写入etcd
	resJson, err := json.Marshal(pod)
	if err != nil {
		log.ErrorLog("UpdatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if pod.Metadata.Namespace == "" || pod.Metadata.Name == "" {
		log.ErrorLog("UpdatePod: namespace or name is empty")
		c.JSON(400, gin.H{"error": "namespace or name is empty"})
		return
	}
	key = config.EtcdPodPrefix + "/" + pod.Metadata.Namespace + "/" + pod.Metadata.Name
	err = etcdclient.EtcdStore.Put(key, string(resJson))
	if err != nil {
		log.ErrorLog("UpdatePod: " + err.Error())
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
		log.ErrorLog("DeletePods: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	key := config.EtcdPodPrefix + "/" + namespace + "/" + name
	log.WarnLog("DeletePods: " + namespace + "/" + name)
	err := etcdclient.EtcdStore.Delete(key)
	if err != nil {
		log.ErrorLog("DeletePods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// TODO 删除kube-controller-manager中的Endpoint
	// TODO 发送删除请求到kubelet
	url := config.KubeletLocalURLPrefix + ":" + fmt.Sprint(config.KubeletAPIPort)
	delUri := url + config.PodsURI
	delUri = strings.Replace(delUri, config.NameSpaceReplace, namespace, -1)
	delUri = strings.Replace(delUri, config.NameReplace, name, -1)
	_, err = httprequest.DelObjMsg(delUri)
	if err != nil {
		log.ErrorLog("DeletePods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": "delete success"})
}

// GetPodEphemeralContainers 获取指定Pod的临时容器
func GetPodEphemeralContainers(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")

	log.DebugLog("GetPodEphemeralContainers: " + namespace + "/" + name)
}

// UpdatePodEphemeralContainers 更新指定Pod的临时容器
func UpdatePodEphemeralContainers(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")

	log.DebugLog("UpdatePodEphemeralContainers: " + namespace + "/" + name)
}

// GetPodLog 获取指定Pod的日志
func GetPodLog(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")

	log.DebugLog("GetPodLog: " + namespace + "/" + name)
}

// GetPodStatus 获取指定Pod的状态
func GetPodStatus(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
	} else if name == "" {
		log.ErrorLog("GetPodStatus: name is empty")
		c.JSON(400, gin.H{"error": "name is empty"})
		return
	}
	key := config.EtcdPodPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("GetPodStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resJson, err := json.Marshal(res)
	if err != nil {
		log.ErrorLog("GetPodStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	pod := &apiObject.Pod{}
	err = json.Unmarshal(resJson, pod)
	if err != nil {
		log.ErrorLog("GetPodStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	//获取pod状态
	status := pod.Status
	byteStatus, err := json.Marshal(status)
	if err != nil {
		log.ErrorLog("GetPodStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": byteStatus})

	log.DebugLog("GetPodStatus: " + namespace + "/" + name)
}

// UpdatePodStatus 更新指定Pod的状态
func UpdatePodStatus(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")

	log.DebugLog("UpdatePodStatus: " + namespace + "/" + name)
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
		log.ErrorLog("GetPods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resJson, err := json.Marshal(res)
	if err != nil {
		log.ErrorLog("GetPods: " + err.Error())
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
		log.ErrorLog("CreatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	newPodName := pod.Metadata.Name
	newPodNamespace := pod.Metadata.Namespace
	if newPodName == "" || newPodNamespace == "" {
		log.ErrorLog("CreatePod: name or namespace is empty")
		c.JSON(400, gin.H{"error": "name or namespace is empty"})
		return
	}
	key := config.EtcdPodPrefix + "/" + pod.Metadata.Namespace + "/" + pod.Metadata.Name
	response, _ := etcdclient.EtcdStore.Get(key)
	if response != "" {
		log.ErrorLog("CreatePod: Pod already exists" + response)
		c.JSON(400, gin.H{"error": "Pod already exists"})
		return
	}
	log.InfoLog("CreatePod: " + newPodNamespace + "/" + newPodName)
	// TODO: 生成 UUID
	pod.Metadata.UUID = uuid.New().String()
	// TODO: 发送的时候筛选 node
	// SchedulerUrl := config.SchedulerURL() + config.SchedulerConfigPath
	// resp, err := http.Get(SchedulerUrl)
	// if err != nil {
	// 	log.ErrorLog("CreatePod: " + err.Error())
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }
	var node apiObject.Node
	res, err := etcdclient.EtcdStore.PrefixGet(config.EtcdNodePrefix)
	if err != nil {
		log.WarnLog("GetNodes: " + err.Error())
		c.JSON(config.HttpErrorCode, gin.H{"error": err.Error()})
		return
	}
	err = json.Unmarshal([]byte(res[0]), &node)
	// err = json.NewDecoder(resp.Body).Decode(&node)
	// if err != nil {
	// 	log.ErrorLog("CreatePod: " + err.Error())
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }
	pod.Spec.NodeName = node.Metadata.Name
	url := config.KubeletLocalURLPrefix + ":" + fmt.Sprint(config.KubeletAPIPort)
	createUri := url + config.PodsURI
	createUri = strings.Replace(createUri, config.NameSpaceReplace, newPodNamespace, -1)
	createUri = strings.Replace(createUri, config.NameReplace, newPodName, -1)
	fmt.Println("createUri: ", createUri)
	reaJson, err := json.Marshal(pod)
	if err != nil {
		log.ErrorLog("CreatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(key, string(reaJson))
	if err != nil {
		log.ErrorLog("CreatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resp, err = httprequest.PostObjMsg(createUri, pod)
	if err != nil {
		log.ErrorLog("Could not post the object message.\n" + err.Error())
		os.Exit(1)
	}
	c.JSON(200, gin.H{"data": resp})
}

// DeletePods 删除所有Pod
func DeletePods(c *gin.Context) {
	namespace := c.Param("namespace")

	log.DebugLog("DeletePods: " + namespace)
}

// GetGlobalPods 获取全局所有Pod
func GetGlobalPods(c *gin.Context) {
	key := config.EtcdPodPrefix
	res, err := etcdclient.EtcdStore.PrefixGet(key)
	if err != nil {
		log.ErrorLog("GetGlobalPods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	resJson, err := json.Marshal(res)
	if err != nil {
		log.ErrorLog("GetGlobalPods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": resJson})
}

// UpdatePodProps TODO: 更新Pod
func UpdatePodProps(new *apiObject.Pod) {
	podBytes, err := json.Marshal(new)
	if err != nil {
		log.ErrorLog("UpdatePodProps: " + err.Error())
		return
	}
	updateUri := config.KubeletLocalURLPrefix + ":" + fmt.Sprint(config.KubeletAPIPort) + config.PodURI
	updateUri = strings.Replace(updateUri, config.NameSpaceReplace, new.Metadata.Namespace, -1)
	updateUri = strings.Replace(updateUri, config.NameReplace, new.Metadata.Name, -1)
	resp, err := httprequest.PutObjMsg(updateUri, podBytes)
	if err != nil {
		log.ErrorLog("UpdatePodProps: " + err.Error())
	}
	//将resp中更新的pod信息存入new
	new = &apiObject.Pod{}
	err = json.NewDecoder(resp.Body).Decode(new)
	if err != nil {
		log.ErrorLog("UpdatePodProps: " + err.Error())
	}
}
