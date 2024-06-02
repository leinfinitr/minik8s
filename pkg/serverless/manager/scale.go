package manager

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"

	httprequest "minik8s/tools/httpRequest"
)

type ScaleManagerImpl struct {
	// 扩容阈值，当同时处理的请求超过这个值时，自动扩容
	Threshold int

	// 每个Serverless Function的实例数量
	FunctionInstanceNum map[string]int
	// 每个Serverless Function的请求数量
	FunctionRequestNum map[string]int

	// 所有运行实例
	Instance map[string]apiObject.Pod
	// 每个实例当前处理的请求数量
	InstanceRequestNum map[string]int
	// 每个实例距离最后一次处理请求所经过的周期
	InstanceLastRequestTime map[string]int

	// 每个Serverless Function对应的Pod
	Pod map[string]apiObject.Pod
	// 每个Serverless Function对应的Serverless
	Serverless map[string]apiObject.Serverless
}

var ScaleManager *ScaleManagerImpl = nil

// NewScaleManager 创建一个新的ScaleManager
func NewScaleManager() *ScaleManagerImpl {
	if ScaleManager == nil {
		ScaleManager = &ScaleManagerImpl{
			Threshold:               10,
			FunctionInstanceNum:     make(map[string]int),
			FunctionRequestNum:      make(map[string]int),
			Instance:                make(map[string]apiObject.Pod),
			InstanceRequestNum:      make(map[string]int),
			InstanceLastRequestTime: make(map[string]int),
			Pod:                     make(map[string]apiObject.Pod),
			Serverless:              make(map[string]apiObject.Serverless),
		}

	}
	return ScaleManager
}

// Run 启动自动扩容控制
func (s *ScaleManagerImpl) Run() {
	// 定时循环检查每个Serverless Function的请求数量和实例数量，根据阈值自动扩容或缩容
	go func() {
		for {
			for name, requestNum := range s.FunctionRequestNum {
				// 扩容
				if requestNum > s.Threshold*s.FunctionInstanceNum[name] {
					s.IncreaseInstance(name)
					continue
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			for name, LastRequestTime := range s.InstanceLastRequestTime {
				// 缩容
				s.InstanceLastRequestTime[name]++
				if LastRequestTime > 30000 {
					s.DecreaseInstance(name)
					continue
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

// IncreaseRequestNum 增加一个Serverless Function的请求数量
func (s *ScaleManagerImpl) IncreaseRequestNum(name string) {
	s.FunctionRequestNum[name]++
}

// DecreaseRequestNum 减少一个Serverless Function的请求数量
func (s *ScaleManagerImpl) DecreaseRequestNum(name string) {
	s.FunctionRequestNum[name]--
}

// IncreaseInstance 增加一个Serverless Function的实例
//
//	name: Serverless Function的名字
func (s *ScaleManagerImpl) IncreaseInstance(name string) {
	// 修改 pod 的 name 和 container name 为 name-InstanceNum
	pod := s.Pod[name]
	instanceName := name + "-" + fmt.Sprint(s.FunctionInstanceNum[name])
	pod.Metadata.Name = instanceName
	pod.Spec.Containers[0].Name = instanceName
	// 转发给 apiServer 创建一个 Pod
	url := config.APIServerURL() + config.PodsURI
	url = strings.Replace(url, config.NameSpaceReplace, pod.Metadata.Namespace, -1)
	res, err := httprequest.PostObjMsg(url, pod)
	if err != nil {
		log.ErrorLog("Could not post the object message." + err.Error())
		os.Exit(1)
	}
	if res.StatusCode != 201 {
		log.ErrorLog("Could not create " + name)
		os.Exit(1)
	}
	// 添加到 Instance 中
	s.FunctionInstanceNum[name]++
	s.Instance[instanceName] = pod
	s.InstanceRequestNum[instanceName] = 0
	s.InstanceLastRequestTime[instanceName] = 0

	log.InfoLog("Create a new pod for " + name + " with name " + instanceName)
}

// DecreaseInstance 删除一个Serverless Function的实例
//
//	instanceName: 运行实例的名字
func (s *ScaleManagerImpl) DecreaseInstance(instanceName string) {
	podInstance := s.Instance[instanceName]
	podName := strings.Split(instanceName, "-")[0]
	// 转发给 apiServer 删除一个 Pod
	url := config.APIServerURL() + config.PodURI
	url = strings.Replace(url, config.NameSpaceReplace, podInstance.Metadata.Namespace, -1)
	url = strings.Replace(url, config.NameReplace, podInstance.Metadata.Name, -1)
	_, err := httprequest.DelMsg(url, nil)
	if err != nil {
		log.ErrorLog("Could not delete the object message." + err.Error())
		os.Exit(1)
	}
	// 从 Instance 中删除
	s.FunctionInstanceNum[podName]--
	delete(s.Instance, instanceName)
	delete(s.InstanceRequestNum, instanceName)
	delete(s.InstanceLastRequestTime, instanceName)

	log.InfoLog("Delete pod " + instanceName + " for " + podName)
}

// RunFunction 运行Serverless Function
//
//	name: Serverless Function的名字
func (s *ScaleManagerImpl) RunFunction(name string, param string) string {
	// 如果当前实例数量不大于0，则循环等待
	for !(s.FunctionInstanceNum[name] > 0) {
		time.Sleep(1 * time.Second)
	}
	// 遍历所有实例，找到一个属于当前Function且请求最少的实例
	minRequestNum := math.MaxInt
	minRequestInstanceName := ""
	for instanceName, requestNum := range s.InstanceRequestNum {
		if strings.HasPrefix(instanceName, name) && requestNum < minRequestNum {
			minRequestNum = s.InstanceRequestNum[instanceName]
			minRequestInstanceName = instanceName
		}
	}
	log.DebugLog("Run function " + name + " with param " + param + " on " + minRequestInstanceName)
	// 如果没有找到合适的实例，则直接报错
	if minRequestInstanceName == "" {
		log.ErrorLog("Could not find a suitable instance for " + name)
		os.Exit(1)
	}
	// 取出该实例
	pod := s.Instance[minRequestInstanceName]
	serverless := s.Serverless[name]
	// 增加该实例处理的请求数量
	s.InstanceRequestNum[minRequestInstanceName]++
	// 重置该实例的最后一次处理请求所经过的周期
	s.InstanceLastRequestTime[minRequestInstanceName] = 0
	// 转发给 apiServer 运行 Pod 中的容器
	url := config.APIServerURL() + config.PodExecURI
	url = strings.Replace(url, config.NameSpaceReplace, pod.Metadata.Namespace, -1)
	url = strings.Replace(url, config.NameReplace, pod.Metadata.Name, -1)
	url = strings.Replace(url, config.ContainerReplace, pod.Spec.Containers[0].Name, -1)
	url = strings.Replace(url, config.ParamReplace, serverless.Command+" "+param, -1)
	response, err := httprequest.GetMsg(url)
	if err != nil {
		log.ErrorLog("Could not post the message." + err.Error())
		os.Exit(1)
	}
	// 减少该实例处理的请求数量
	s.InstanceRequestNum[minRequestInstanceName]--
	// 返回结果
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.ErrorLog(err.Error())
		os.Exit(1)
	}
	result := string(body)
	// 去掉result中首尾的引号
	result = result[1 : len(result)-1]
	log.InfoLog("Run function " + name + " with param " + param + " result: " + result)
	return result
}

// AddPod 添加一个Pod
func (s *ScaleManagerImpl) AddPod(pod apiObject.Pod) {
	log.DebugLog("Add pod " + pod.Metadata.Name)
	s.Pod[pod.Metadata.Name] = pod
	s.FunctionInstanceNum[pod.Metadata.Name] = 0
	s.FunctionRequestNum[pod.Metadata.Name] = 0
}

// DeletePod 删除一个Pod
func (s *ScaleManagerImpl) DeletePod(name string) {
	log.DebugLog("Delete pod " + name)
	delete(s.Pod, name)
	delete(s.FunctionInstanceNum, name)
	delete(s.FunctionRequestNum, name)
}

// AddServerless 添加一个Serverless
func (s *ScaleManagerImpl) AddServerless(serverless apiObject.Serverless) {
	s.Serverless[serverless.Name] = serverless
}

// DeleteServerless 删除一个Serverless
func (s *ScaleManagerImpl) DeleteServerless(name string) {
	log.DebugLog("Delete serverless " + name)
	delete(s.Serverless, name)
	// 删除该 Serverless 对应的所有实例
	for instanceName := range s.Instance {
		if strings.HasPrefix(instanceName, name) {
			s.DecreaseInstance(instanceName)
		}
	}
}

// GetAllServerless 获取所有的Serverless
func (s *ScaleManagerImpl) GetAllServerless() []apiObject.Serverless {
	serverless := make([]apiObject.Serverless, 0)
	for _, s := range s.Serverless {
		serverless = append(serverless, s)
	}
	return serverless
}
