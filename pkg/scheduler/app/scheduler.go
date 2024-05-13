package scheduler

import (
	"encoding/json"
	"fmt"
	httprequest "minik8s/internal/pkg/httpRequest"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"
	"sync"
)

type Scheduler struct {
	// ApiServerConfig 存储apiServer的配置信息，用于和apiServer进行通信
	ApiServerConfig *config.APIServerConfig
	//调度策略
	Policy string
}

const (
	RoundRobin = "RoundRobin"
)

var glbcnt int
var lock sync.Mutex

func (s *Scheduler) schedRequest() string {
	// 从apiServer获取pod信息
	podList := s.getNodesList()
	// 调度pod
	return s.schedule(podList)
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		ApiServerConfig: config.NewAPIServerConfig(),
		Policy:          RoundRobin,
	}
}

func (s *Scheduler) getNodesList() []apiObject.Pod {
	// 从apiServer获取所有的pod信息
	url := s.ApiServerConfig.APIServerURL() + config.NodesURI
	var podList []apiObject.Pod
	resp, err := httprequest.GetObjMsg(url, &podList, "data")
	if err != nil {
		log.DebugLog("httprequest.GetObjMsg err:" + err.Error())
		return nil
	}
	if resp.StatusCode != 200 {
		log.DebugLog("httprequest.GetObjMsg StatusCode:" + fmt.Sprint(resp.StatusCode))
		return nil
	}
	return podList
}

func (s *Scheduler) schedule(podList []apiObject.Pod) string {
	switch s.Policy {
	case RoundRobin:
		return s.roundRobinSched(podList)
	default:
		return s.roundRobinSched(podList)
	}
}

func (s *Scheduler) roundRobinSched(podList []apiObject.Pod) string {
	lock.Lock()
	defer lock.Unlock()
	if glbcnt >= len(podList) {
		glbcnt = 0
	}
	pod := podList[glbcnt]
	glbcnt++
	data, err := json.Marshal(pod)
	if err != nil {
		log.DebugLog("json.Marshal err:" + err.Error())
		return ""
	}
	return string(data)
}

func Run() string {
	scheduler := NewScheduler()
	data := scheduler.schedRequest()
	if data == "" {
		log.DebugLog("schedRequest data is nil")
	}
	log.DebugLog("schedRequest data:" + data)
	return data
}
