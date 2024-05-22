package scale

import (
	"minik8s/pkg/apiObject"
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
}

var ScaleManager *ScaleManagerImpl = nil

// NewScaleManager 创建一个新的ScaleManager
func NewScaleManager() *ScaleManagerImpl {
	if ScaleManager == nil {
		ScaleManager = &ScaleManagerImpl{
			Threshold:   10,
			InstanceNum: make(map[string]int),
			RequestNum:  make(map[string]int),
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
			if requestNum < s.Threshold*(s.InstanceNum[name]-1) {
				s.DecreaseInstanceNum(name)
				continue
			}
		}
		time.Sleep(10 * time.Second)
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
	s.InstanceNum[name]++
	// TODO: 创建一个新的Serverless Function实例
}

// DecreaseInstanceNum 减少一个Serverless Function的实例
func (s *ScaleManagerImpl) DecreaseInstanceNum(name string) {
	s.InstanceNum[name]--
	// TODO: 删除一个Serverless Function实例
}

// RunFunction 运行Serverless Function
func (s *ScaleManagerImpl) RunFunction(pod apiObject.Pod, param string) string {
	// TODO: 取出一个实例运行
	return ""
}
