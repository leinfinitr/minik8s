package scheduler

import (
	"encoding/json"
	"fmt"
	"io"

	// httprequest "minik8s/internal/pkg/httpRequest"
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

var glbcnt int
var lock sync.Mutex

func (s *Scheduler) schedRequest() string {
	// 从apiServer获取pod信息
	podList := s.getNodesList()
	// 调度pod
	return s.schedule(podList)
}

func NewScheduler() *Scheduler {
	glbcnt = 0
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
        log.ErrorLog("httprequest.GetObjMsg err:" + err.Error())
        return nil
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        log.ErrorLog("httprequest.GetObjMsg StatusCode:" + fmt.Sprint(resp.StatusCode))
        return nil
    }
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.ErrorLog("ioutil.ReadAll err:" + err.Error())
		return nil
	}
	log.InfoLog("getNodesList body:" + string(bodyBytes))
    bodyString := string(bodyBytes)
	err = json.Unmarshal([]byte(bodyString), &NodeList)
	if err != nil {
		log.ErrorLog("json.Unmarshal err:" + err.Error())
		return nil
	}
	fmt.Println(NodeList)
    return NodeList
}

func (s *Scheduler) schedule(nodeList []apiObject.Node) string {
	switch s.Policy {
	case RoundRobin:
		return s.roundRobinSched(nodeList)
	default:
		return s.roundRobinSched(nodeList)
	}
}

func (s *Scheduler) roundRobinSched(nodeList []apiObject.Node) string {
	lock.Lock()
	defer lock.Unlock()
	if glbcnt >= len(nodeList) {
		glbcnt = 0
	}
	node := nodeList[glbcnt]
	glbcnt++
	data, err := json.Marshal(node)
	if err != nil {
		log.ErrorLog("json.Marshal err:" + err.Error())
		return ""
	}
	return string(data)
}

func Run() {
	scheduler := NewScheduler()
	r := gin.Default()

	r.GET(config.SchedulerPath(), func(c *gin.Context) {
		data := scheduler.schedRequest()
		c.JSON(200, gin.H{"data": data})
	})

	log.InfoLog("Starting scheduler HTTP server on :7820")
	r.Run(":"+config.SchedulerPort())
}