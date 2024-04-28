package status

import (
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/pkg/kubelet/runtime"
	"minik8s/tools/host"
	"minik8s/tools/log"
	"minik8s/tools/netRequest"
	"strconv"
)

type StatusManager interface {
	// RegisterNode 注册节点
	RegisterNode() error
}

type statusManagerImpl struct {
	runtimeMgr   *runtime.RuntimeManager
	apiServerURL string
}

var statusManager *statusManagerImpl = nil

// GetStatusManager 返回的是接口类型
func GetStatusManager(apiServerURL string) StatusManager {
	if statusManager == nil {
		statusManager = &statusManagerImpl{
			runtimeMgr:   runtime.GetRuntimeManager(),
			apiServerURL: apiServerURL,
		}
	}

	return statusManager
}

// RegisterNode 在kubelet刚开始创建时，需要到apiServer的work node去注册
//
//	通过发送POST请求的方式去注册，默认API："/api/v1/nodes"
func (s *statusManagerImpl) RegisterNode() error {
	// 注册所需的参数
	HostName, err := host.GetHostname()
	HostIP, err := host.GetHostIP()

	// 获取主机的内存大小
	capacity := make(map[string]string)
	totalMemory, err := host.GetTotalMemory()
	capacity["memory"] = strconv.FormatUint(totalMemory, 10)

	// 获取主机的内存和CPU使用率
	allocatable := make(map[string]string)
	MemoryUsage, err := host.GetMemoryUsageRate()
	CPUUsage, err := host.GetCPULoad()
	allocatable["memory"] = strconv.FormatFloat(MemoryUsage, 'f', -1, 64)
	allocatable["cpu"] = strconv.FormatFloat(CPUUsage[0], 'f', -1, 64)

	node := &apiObject.Node{
		TypeMeta: apiObject.TypeMeta{
			Kind:       "Node",
			APIVersion: "v1",
		},
		Metadata: apiObject.ObjectMeta{
			Name:        HostName,
			Namespace:   "", // 该字段为空
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			UUID:        "", //	由API Server生成
		},
		Spec: apiObject.NodeSpec{
			PodCIDR:       "", // 未使用
			ProviderID:    "", // 未使用
			Unschedulable: true,
		},
		Status: apiObject.NodeStatus{
			Capacity:    capacity,
			Allocatable: allocatable,
			Phase:       "running",
			Conditions: []apiObject.NodeCondition{
				{
					Type:   "Ready", // Ready: kubelet准备好接受Pod
					Status: "True",
				},
			},
			Addresses: []apiObject.NodeAddress{
				{
					Type:    "InternalIP",
					Address: HostIP,
				},
			},
		},
	}

	url := "http://" + s.apiServerURL + "/api/v1/nodes"

	statusCode, _, err := netRequest.PostRequestByTarget(url, node)
	if err != nil {
		return err
	}

	if statusCode != config.HttpSuccessCode {
		log.ErrorLog("register node failed")
	}
	log.InfoLog("register node success")

	return nil
}
