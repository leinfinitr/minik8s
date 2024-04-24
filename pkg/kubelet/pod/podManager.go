package pod

import (
	"minik8s/pkg/apiObject"
	"minik8s/pkg/kubelet/runtime"
)

type PodUUID string
type EventType string

/*  */
type PodManager struct {
	/* 实现从UUID到pod的映射 */
	podMapByUUID map[PodUUID]*apiObject.Pod

	eventQueue chan EventType

	/* 不同事件的处理函数 */

	addPodHandler            func(pod *apiObject.Pod) error
	startPodHandler          func(pod *apiObject.Pod) error
	restartPodHandler        func(pod *apiObject.Pod) error
	stopPodHandler           func(pod *apiObject.Pod) error
	deletePodHandler         func(pod *apiObject.Pod) error
	recreateContainerHandler func(pod *apiObject.Pod) error
	execPodHandler           func(pod *apiObject.Pod) error
}

/* Singleton pattern */
var podManager *PodManager = nil

func GetPodManager() *PodManager {
	if podManager == nil {
		newMapUUIDToPod := make(map[PodUUID]*apiObject.Pod)
		eventChan := make(chan EventType)
		// TODO：此处需要获取所有pod的信息，接口应当放在podUtils中，来更新map，未实现
		runtimeMgr := runtime.GetRuntimeManager()
		podManager = &PodManager{
			podMapByUUID:             newMapUUIDToPod,
			eventQueue:               eventChan,
			addPodHandler:            runtimeMgr.CreatePod,
			startPodHandler:          runtimeMgr.StartPod,
			restartPodHandler:        runtimeMgr.RestartPod,
			stopPodHandler:           runtimeMgr.StopPod,
			deletePodHandler:         runtimeMgr.DeletePod,
			recreateContainerHandler: runtimeMgr.RecreatePodContainer,
			execPodHandler:           runtimeMgr.ExecPodContainer,
		}
	}

	return podManager
}
