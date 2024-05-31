package scheduler

import (
	"encoding/json"
	"fmt"
	"io"

	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
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

var globalCount int
var lock sync.Mutex

func (s *Scheduler) scheduleRequest() apiObject.Node {
	// 从apiServer获取pod信息
	podList := s.getNodesList()
	// 调度pod
	return s.schedule(podList)
}

func NewScheduler() *Scheduler {
	globalCount = 0
	return &Scheduler{
		ApiServerConfig: config.NewAPIServerConfig(),
		Policy:          RoundRobin,
	}
}

func (s *Scheduler) getNodesList() []apiObject.Node {
	// 从apiServer获取所有的pod信息
	url := s.ApiServerConfig.APIServerURL() + config.NodesURI
	var NodeList []apiObject.Node
	resp, err := http.Get(url)
	if err != nil {
		log.ErrorLog("http request.GetObjMsg err:" + err.Error())
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.ErrorLog("http request.GetObjMsg StatusCode:" + fmt.Sprint(resp.StatusCode))
		return nil
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.ErrorLog("io.ReadAll err:" + err.Error())
		return nil
	}
	log.DebugLog("getNodesList body:" + string(bodyBytes))
	bodyString := string(bodyBytes)
	err = json.Unmarshal([]byte(bodyString), &NodeList)
	if err != nil {
		log.ErrorLog("json.Unmarshal err:" + err.Error())
		return nil
	}
	fmt.Println(NodeList)
	return NodeList
}

func (s *Scheduler) schedule(nodeList []apiObject.Node) apiObject.Node {
	switch s.Policy {
	case RoundRobin:
		return s.roundRobin(nodeList)
	default:
		return s.roundRobin(nodeList)
	}
}

func (s *Scheduler) roundRobin(nodeList []apiObject.Node) apiObject.Node {
	lock.Lock()
	defer lock.Unlock()
	if globalCount >= len(nodeList) {
		globalCount = 0
	}
	node := nodeList[globalCount]
	globalCount++
	return node
}

func Run() {
	gin.SetMode(gin.ReleaseMode)
	scheduler := NewScheduler()
	r := gin.New()

	r.GET(config.SchedulerPath(), func(c *gin.Context) {
		data := scheduler.scheduleRequest()
		c.JSON(200, data)
	})

	log.InfoLog("Starting scheduler HTTP server on :7820")
	r.Run(":" + config.SchedulerPort())
}
