package status

import (
	"minik8s/pkg/kubelet/runtime"
	// "minik8s/pkg/apiObject"
)

/* statusManager 功能介绍
 * 1. 是用来记录pod各种状态信息的模块
 * 2. 该模块负责提供获取相关信息的接口，会有其他组件通过接口来监控pod的状态
 * 3. 定时发送pod的状态消息给apiServer，该信息同样会被当作心跳
 */

type StatusManager interface {
	// GerPodInfoFromCRI()

	// 注册节点
	RegisterNode() error
}

type statusManagerImpl struct {
	runtimeMgr   *runtime.RuntimeManager
	apiServerIP  string
	apiServerURL string
}

var statusManager *statusManagerImpl = nil

// 返回的是接口类型
func GetStatusManager(apiServerURL string, apiServerIP string) (StatusManager, error) {
	if statusManager == nil {
		statusManager = &statusManagerImpl{
			runtimeMgr:   runtime.GetRuntimeManager(),
			apiServerIP:  apiServerIP,
			apiServerURL: apiServerURL,
		}
	}

	return statusManager, nil
}

/* 在kubelet刚开始创建时，需要到apiServer的work node去注册
 * 通过发送POST请求的方式去注册
 * 默认API："/api/v1/nodes"
 */
func (s *statusManagerImpl) RegisterNode() error {
	// TODO: check if this node has been registered

	// node := &apiObject.Node{

	// }
	return nil
}
