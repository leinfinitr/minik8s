package scale

import (
	"fmt"
	"io"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	httprequest "minik8s/tools/httpRequest"
	"minik8s/tools/log"
	"os"
	"strings"
	"time"
)

type ScaleManagerImpl struct {
	// 扩容阈值，当同时处理的请求超过这个值时，自动扩容
	Threshold int
	// 当前每个Serverless Function的实例数量
	InstanceNum map[string]int
	// 当前每个Serverless Function的请求数量
	RequestNum map[string]int
	// 每个Serverless Function的所有运行实例
	Instance map[string][]apiObject.Pod
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
			Threshold:   10,
			InstanceNum: make(map[string]int),
			RequestNum:  make(map[string]int),
			Instance:    make(map[string][]apiObject.Pod),
			Pod:         make(map[string]apiObject.Pod),
			Serverless:  make(map[string]apiObject.Serverless),
		}

	}
	return ScaleManager
}

// Run 启动自动扩容控制
func (s *ScaleManagerImpl) Run() {
	// 定时循环检查每个Serverless Function的请求数量和实例数量，根据阈值自动扩容或缩容
	for {
		for name, requestNum := range s.RequestNum {
			// 扩容
			if requestNum > s.Threshold*s.InstanceNum[name] {
				s.IncreaseInstanceNum(name)
				continue
			}
			if s.InstanceNum[name] == 0 {
				continue
			}
			// 缩容
			if requestNum <= s.Threshold*(s.InstanceNum[name]-1) {
				s.DecreaseInstanceNum(name)
				continue
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// IncreaseRequestNum 增加一个Serverless Function的请求数量
func (s *ScaleManagerImpl) IncreaseRequestNum(name string) {
	s.RequestNum[name]++
}

// DecreaseRequestNum 减少一个Serverless Function的请求数量
func (s *ScaleManagerImpl) DecreaseRequestNum(name string) {
	s.RequestNum[name]--
}

// IncreaseInstanceNum 增加一个Serverless Function的实例
func (s *ScaleManagerImpl) IncreaseInstanceNum(name string) {
	// 修改 pod 的 name 和 container name 为 name-InstanceNum
	pod := s.Pod[name]
	pod.Metadata.Name = name + "-" + fmt.Sprint(s.InstanceNum[name])
	pod.Spec.Containers[0].Name = name + "-" + fmt.Sprint(s.InstanceNum[name])
	// 转发给 apiServer 创建一个 Pod
	url := config.APIServerURL() + config.PodsURI
	url = strings.Replace(url, config.NameSpaceReplace, pod.Metadata.Namespace, -1)
	_, err := httprequest.PostObjMsg(url, pod)
	if err != nil {
		log.ErrorLog("Could not post the object message." + err.Error())
		os.Exit(1)
	}
	// 添加到 Instance 中
	s.Instance[name] = append(s.Instance[name], pod)
	s.InstanceNum[name]++

	log.InfoLog("Create a new pod for " + name + " with name " + pod.Metadata.Name)
}

// DecreaseInstanceNum 减少一个Serverless Function的实例
func (s *ScaleManagerImpl) DecreaseInstanceNum(name string) {
	s.InstanceNum[name]--
	// TODO: 删除一个Serverless Function实例
}

// RunFunction 运行Serverless Function
func (s *ScaleManagerImpl) RunFunction(name string, param string) string {
	// 如果当前没有实例，则循环等待
	for s.InstanceNum[name] == 0 {
		time.Sleep(1 * time.Second)
	}
	// 从 Instance 中取出最后一个 Pod
	pod := s.Instance[name][s.InstanceNum[name]-1]
	serverless := s.Serverless[name]
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
	// 返回结果
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.ErrorLog("Could not read the response body." + err.Error())
		os.Exit(1)
	}
	return string(body)
}

// AddPod 添加一个Pod
func (s *ScaleManagerImpl) AddPod(pod apiObject.Pod) {
	s.Pod[pod.Metadata.Name] = pod
}

// AddServerless 添加一个Serverless
func (s *ScaleManagerImpl) AddServerless(serverless apiObject.Serverless) {
	s.Serverless[serverless.Name] = serverless
}
