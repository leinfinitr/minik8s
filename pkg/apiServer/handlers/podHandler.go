package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/pkg/controller"
	"minik8s/tools/log"

	etcdclient "minik8s/pkg/apiServer/etcdClient"
	httprequest "minik8s/tools/httpRequest"
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
	rep, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("UpdatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if rep == "" {
		log.ErrorLog("UpdatePod: pod not found")
		c.JSON(400, gin.H{"error": "pod not found"})
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
	// 解析更新后的pod
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
	// 将更新后的pod写入etcd
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
	log.InfoLog("DeletePods: " + namespace + "/" + name)

	key := config.EtcdPodPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("DeletePods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	pod := &apiObject.Pod{}
	err = json.Unmarshal([]byte(res), pod)
	if err != nil {
		log.ErrorLog("DeletePods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	nodeName := pod.Spec.NodeName
	// 删除etcd中的Pod
	err = etcdclient.EtcdStore.Delete(key)
	if err != nil {
		log.ErrorLog("DeletePods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 若pod使用了pvc，则对其进行处理
	if pod.Spec.Volumes != nil {
		for _, volume := range pod.Spec.Volumes {
			if volume.PersistentVolumeClaim != nil {
				pvcKey := config.EtcdPvcPrefix + "/" + pod.Metadata.Namespace + "/" + volume.PersistentVolumeClaim.ClaimName
				pvcResponse, _ := etcdclient.EtcdStore.Get(pvcKey)
				// 若pvc不存在，则跳过
				if pvcResponse == "" {
					continue
				}
				// 若pvc存在，则将其与pod解绑
				pvc := &apiObject.PersistentVolumeClaim{}
				err = json.Unmarshal([]byte(pvcResponse), pvc)
				if err != nil {
					log.ErrorLog("DeletePods: " + err.Error())
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				err = controller.ControllerManagerInstance.UnbindPodToPvc(pvc)
				if err != nil {
					log.ErrorLog("DeletePods: " + err.Error())
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
			}
		}
	}
	// 获取pod所在node的IP
	res, err = etcdclient.EtcdStore.Get(config.EtcdNodePrefix + "/" + nodeName)
	if err != nil {
		log.ErrorLog("DeletePods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	node := &apiObject.Node{}
	err = json.Unmarshal([]byte(res), node)
	if err != nil {
		log.ErrorLog("DeletePods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	addresses := node.Status.Addresses
	address := addresses[0].Address
	// 发送删除请求到kubelet
	url := config.HttpSchema + address + ":" + fmt.Sprint(config.KubeletAPIPort)
	delUri := url + config.PodURI
	delUri = strings.Replace(delUri, config.NameSpaceReplace, namespace, -1)
	delUri = strings.Replace(delUri, config.NameReplace, name, -1)
	_, err = httprequest.DelMsg(delUri, *pod)
	if err != nil {
		log.ErrorLog("DeletePods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": "delete success"})
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

	log.DebugLog("GetPodStatus: " + namespace + "/" + name)
	c.JSON(200, gin.H{"data": byteStatus})
}

// UpdatePodStatus 更新指定Pod的状态
func UpdatePodStatus(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	if name == "" || namespace == "" {
		log.ErrorLog("UpdatePodStatus: name or namespace is empty")
		c.JSON(400, gin.H{"error": "name or namespace is empty"})
		return
	}
	log.InfoLog("UpdatePodStatus: " + namespace + "/" + name)

	key := config.EtcdPodPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("UpdatePodStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	podStatus := &apiObject.PodStatus{}
	err = c.ShouldBindJSON(podStatus)
	if err != nil {
		log.ErrorLog("UpdatePodStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 更新podStatus
	pod := &apiObject.Pod{}
	err = json.Unmarshal([]byte(res), pod)
	if err != nil {
		log.ErrorLog("UpdatePodStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	pod.Status = *podStatus
	// 将更新后的pod写入etcd
	resJson, err := json.Marshal(pod)
	if err != nil {
		log.ErrorLog("UpdatePodStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = etcdclient.EtcdStore.Put(key, string(resJson))
	if err != nil {
		log.ErrorLog("UpdatePodStatus: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	log.InfoLog("UpdatePodStatus: " + namespace + "/" + name)
	c.JSON(200, gin.H{"data": resJson})
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

	var pods []apiObject.Pod
	for _, v := range res {
		pod := apiObject.Pod{}
		err = json.Unmarshal([]byte(v), &pod)
		if err != nil {
			log.ErrorLog("GetPods: " + err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		pods = append(pods, pod)
	}

	c.JSON(200, pods)
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
	// 检查pod的name和namespace是否为空
	newPodName := pod.Metadata.Name
	newPodNamespace := pod.Metadata.Namespace
	if newPodName == "" || newPodNamespace == "" {
		log.ErrorLog("CreatePod: name or namespace is empty")
		c.JSON(400, gin.H{"error": "name or namespace is empty"})
		return
	}
	// 判断pod是否已经存在
	key := config.EtcdPodPrefix + "/" + pod.Metadata.Namespace + "/" + pod.Metadata.Name
	response, _ := etcdclient.EtcdStore.Get(key)
	if response != "" {
		log.ErrorLog("CreatePod: Pod already exists" + response)
		c.JSON(400, gin.H{"error": "Pod already exists"})
		return
	}
	// 若pod使用了pvc，则对其进行处理
	if pod.Spec.Volumes != nil {
		for _, volume := range pod.Spec.Volumes {
			if volume.PersistentVolumeClaim != nil {
				pvcKey := config.EtcdPvcPrefix + "/" + pod.Metadata.Namespace + "/" + volume.PersistentVolumeClaim.ClaimName
				pvcResponse, _ := etcdclient.EtcdStore.Get(pvcKey)
				// 若pvc不存在，则返回错误
				if pvcResponse == "" {
					log.ErrorLog("CreatePod: PVC not found")
					c.JSON(400, gin.H{"error": "PVC not found"})
					return
				}
				// 若pvc存在，则检查pvc的状态是否为Bound、是否已经绑定到pod
				pvc := &apiObject.PersistentVolumeClaim{}
				err = json.Unmarshal([]byte(pvcResponse), pvc)
				if err != nil {
					log.ErrorLog("CreatePod: " + err.Error())
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				if pvc.Status.Phase != apiObject.ClaimBound || pvc.Status.IsBound {
					log.ErrorLog("CreatePod: PVC can't be used")
					c.JSON(400, gin.H{"error": "PVC can't be used"})
					return
				}
				// 将pod绑定到pvcs
				err = controller.ControllerManagerInstance.BindPodToPvc(pvc, pod.Metadata.Namespace+"/"+pod.Metadata.Name)
				if err != nil {
					log.ErrorLog("CreatePod: " + err.Error())
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				// 为pod中所有使用该持久化卷挂载的容器添加Mount
				for _, container := range pod.Spec.Containers {
					for _, volumeMount := range container.VolumeMounts {
						if volumeMount.Name == volume.Name {
							// 获取pv名称
							pvName := controller.ControllerManagerInstance.GetPvcBind(volume.PersistentVolumeClaim.ClaimName)
							if pvName == "" {
								log.ErrorLog("CreatePod: " + "PersistentVolumeClaim not bind")
								c.JSON(400, gin.H{"error": "PersistentVolumeClaim not bind"})
								return
							}
							// 为容器添加Mount
							container.Mounts = append(container.Mounts, &apiObject.Mount{
								HostPath:      config.PVClientPath + "/" + pvName,
								ContainerPath: volumeMount.MountPath,
								ReadOnly:      volume.PersistentVolumeClaim.ReadOnly,
							})
						}
					}
				}
			}
		}
	}
	// 完成检查，生成 UUID
	pod.Metadata.UUID = uuid.New().String()
	log.InfoLog("CreatePod: " + newPodNamespace + "/" + newPodName)
	// 发送的时候筛选 node
	ScheduledUri := config.SchedulerURL() + config.SchedulerConfigPath
	resp, err := http.Get(ScheduledUri)
	if err != nil {
		log.ErrorLog("CreatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	var node apiObject.Node
	err = json.NewDecoder(resp.Body).Decode(&node)
	if err != nil {
		log.ErrorLog("CreatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	pod.Spec.NodeName = node.Metadata.Name
	// 得到node的IP
	addresses := node.Status.Addresses
	address := addresses[0].Address
	log.DebugLog("CreatePod: " + address)
	// 发送创建请求到kubelet
	url := config.HttpSchema + address + ":" + fmt.Sprint(config.KubeletAPIPort)
	createUri := url + config.PodsURI
	createUri = strings.Replace(createUri, config.NameSpaceReplace, newPodNamespace, -1)
	createUri = strings.Replace(createUri, config.NameReplace, newPodName, -1)
	log.DebugLog("createUri: " + createUri)
	// 发送创建请求并解析返回的pod信息
	resp, err = httprequest.PostObjMsg(createUri, pod)
	if err != nil {
		log.ErrorLog("Could not post the object message.\n" + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.ErrorLog("CreatePod: " + resp.Status)
		c.JSON(500, gin.H{"error": resp.Status})
		return
	}
	err = json.NewDecoder(resp.Body).Decode(&pod)
	if err != nil {
		log.ErrorLog("CreatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	reaJson, err := json.Marshal(pod)
	if err != nil {
		log.ErrorLog("CreatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 将pod信息存入etcd
	err = etcdclient.EtcdStore.Put(key, string(reaJson))
	if err != nil {
		log.ErrorLog("CreatePod: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, reaJson)
}

// DeletePods 删除所有Pod
func DeletePods(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default"
		c.Params = append(c.Params, gin.Param{Key: "namespace", Value: namespace})
	}

	key := config.EtcdPodPrefix + "/" + namespace
	res, err := etcdclient.EtcdStore.PrefixGet(key)
	if err != nil {
		log.ErrorLog("DeletePods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// 遍历etcd中的所有pod，调用DeletePod删除
	for _, v := range res {
		pod := &apiObject.Pod{}
		err = json.Unmarshal([]byte(v), pod)
		if err != nil {
			log.ErrorLog("DeletePods: " + err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Params = append(c.Params, gin.Param{Key: "name", Value: pod.Metadata.Name})
		DeletePod(c)
	}

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

	var pods []apiObject.Pod
	for _, v := range res {
		pod := apiObject.Pod{}
		err = json.Unmarshal([]byte(v), &pod)
		if err != nil {
			log.ErrorLog("GetPods: " + err.Error())
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		pods = append(pods, pod)
	}

	c.JSON(200, pods)
}

// UpdatePodProps 更新Pod
func UpdatePodProps(new *apiObject.Pod) {
	podBytes, err := json.Marshal(new)
	if err != nil {
		log.ErrorLog("UpdatePodProps: " + err.Error())
		return
	}
	nodeName := new.Spec.NodeName

	key := config.EtcdNodePrefix + "/" + nodeName
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("UpdatePodProps: " + err.Error())
		return
	}

	node := &apiObject.Node{}
	err = json.Unmarshal([]byte(res), node)
	if err != nil {
		log.ErrorLog("UpdatePodProps: " + err.Error())
		return
	}
	addresses := node.Status.Addresses
	address := addresses[0].Address

	url := config.HttpSchema + address + ":" + fmt.Sprint(config.KubeletAPIPort)
	updateUri := url + config.PodURI
	updateUri = strings.Replace(updateUri, config.NameSpaceReplace, new.Metadata.Namespace, -1)
	updateUri = strings.Replace(updateUri, config.NameReplace, new.Metadata.Name, -1)
	resp, err := httprequest.PutObjMsg(updateUri, podBytes)
	if err != nil {
		log.ErrorLog("UpdatePodProps: " + err.Error())
		return
	}
	// 将resp中更新的pod信息存入new
	new = &apiObject.Pod{}
	err = json.NewDecoder(resp.Body).Decode(new)
	if err != nil {
		log.ErrorLog("UpdatePodProps: " + err.Error())
	}
}

// ExecPod 根据Pod和container名称执行相应的命令
func ExecPod(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	container := c.Param("container")
	param := c.Param("param")
	log.DebugLog("ExecPod: " + namespace + "/" + name + "/" + container)

	// 取出Pod
	key := config.EtcdPodPrefix + "/" + namespace + "/" + name
	res, err := etcdclient.EtcdStore.Get(key)
	if err != nil {
		log.ErrorLog("ExecPod: " + err.Error())
		c.JSON(500, err.Error())
		return
	}
	pod := &apiObject.Pod{}
	err = json.Unmarshal([]byte(res), pod)
	if err != nil {
		log.ErrorLog("ExecPod: " + err.Error())
		c.JSON(500, err.Error())
		return
	}
	// 取出containerID
	containerID := ""
	for _, containers := range pod.Spec.Containers {
		if containers.Name == container {
			containerID = containers.ContainerID
			break
		}
	}
	if containerID == "" {
		log.ErrorLog("ExecPod: containerID is empty")
		c.JSON(400, "containerID is empty")
		return
	}

	// 获取pod所在node的IP
	res, err = etcdclient.EtcdStore.Get(config.EtcdNodePrefix + "/" + pod.Spec.NodeName)
	if err != nil {
		log.ErrorLog("DeletePods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	node := &apiObject.Node{}
	err = json.Unmarshal([]byte(res), node)
	if err != nil {
		log.ErrorLog("DeletePods: " + err.Error())
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	addresses := node.Status.Addresses
	address := addresses[0].Address

	// 执行命令
	log.DebugLog("ExecPod: " + namespace + "/" + name + "/" + containerID + "/" + param)
	execUri := config.HttpSchema + address + ":" + fmt.Sprint(config.KubeletAPIPort) + config.PodExecURI
	execUri = strings.Replace(execUri, config.NameSpaceReplace, namespace, -1)
	execUri = strings.Replace(execUri, config.NameReplace, name, -1)
	execUri = strings.Replace(execUri, config.ContainerReplace, containerID, -1)
	execUri = strings.Replace(execUri, config.ParamReplace, param, -1)
	resp, err := httprequest.GetMsg(execUri)
	if err != nil {
		log.ErrorLog("ExecPod: " + err.Error())
		c.JSON(500, err.Error())
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.ErrorLog("ExecPod: " + err.Error())
		c.JSON(500, err.Error())
		return
	}
	result := string(body)
	// 去掉result中首尾的引号
	result = result[1 : len(result)-1]
	log.DebugLog("ExecPod: " + namespace + "/" + name + "/" + containerID + "/" + param + " success: " + result)
	c.JSON(200, result)
}
