package pod

import (
	"errors"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/kubelet/runtime"
	"minik8s/tools/log"
)

type EventType string

type PodManager interface {
	AddPod(pod *apiObject.Pod) error
	DeletePod(pod *apiObject.Pod) error
	StartPod(pod *apiObject.Pod) error
	StopPod(pod *apiObject.Pod) error
	RestartPod(pod *apiObject.Pod) error
	DeletePodByUUID(pod *apiObject.Pod) error
	RecreatePodContainer(pod *apiObject.Pod) error
	ExecPodContainer(pod *apiObject.Pod) error
}

/*  */
type podManagerImpl struct {
	/* 实现从UUID到pod的映射 */
	PodMapByUUID map[string]*apiObject.Pod

	EventQueue chan EventType

	/* 不同事件的处理函数 */

	AddPodHandler            func(pod *apiObject.Pod) error
	StartPodHandler          func(pod *apiObject.Pod) error
	RestartPodHandler        func(pod *apiObject.Pod) error
	StopPodHandler           func(pod *apiObject.Pod) error
	DeletePodHandler         func(pod *apiObject.Pod) error
	RecreateContainerHandler func(pod *apiObject.Pod) error
	ExecPodHandler           func(pod *apiObject.Pod) error
}

/* Singleton pattern */
var podManager *podManagerImpl = nil

func GetPodManager() PodManager {
	if podManager == nil {
		newMapUUIDToPod := make(map[string]*apiObject.Pod)
		eventChan := make(chan EventType)
		// TODO：此处需要获取所有pod的信息，接口应当放在podUtils中，来更新map，未实现

		runtimeMgr := runtime.GetRuntimeManager()
		podManager = &podManagerImpl{
			PodMapByUUID:             newMapUUIDToPod,
			EventQueue:               eventChan,
			AddPodHandler:            runtimeMgr.CreatePod,
			StartPodHandler:          runtimeMgr.StartPod,
			RestartPodHandler:        runtimeMgr.RestartPod,
			StopPodHandler:           runtimeMgr.StopPod,
			DeletePodHandler:         runtimeMgr.DeletePod,
			RecreateContainerHandler: runtimeMgr.RecreatePodContainers,
			ExecPodHandler:           runtimeMgr.ExecPodContainer,
		}
	}

	return podManager
}

func (p *podManagerImpl) AddPod(pod *apiObject.Pod) error {
	log.DebugLog("AddPod")
	uuid := pod.GetPodUUID()
	if _, ok := p.PodMapByUUID[uuid]; ok {
		// 说明接受过这个请求了
		return errors.New("pod message has been handled")
	}

	go func() {
		err := p.AddPodHandler(pod)
		if err != nil {
			log.ErrorLog("AddPodHandler error: " + err.Error())
		} else {
			p.PodMapByUUID[uuid] = pod
			log.InfoLog("AddPodHandler success")
		}
	}()

	return nil
}

func (*podManagerImpl) DeletePod(pod *apiObject.Pod) error {
	return nil
}

func (*podManagerImpl) StartPod(pod *apiObject.Pod) error {
	return nil
}

func (*podManagerImpl) StopPod(pod *apiObject.Pod) error {
	return nil
}

func (*podManagerImpl) RestartPod(pod *apiObject.Pod) error {
	return nil
}

func (*podManagerImpl) DeletePodByUUID(pod *apiObject.Pod) error {
	return nil
}

func (*podManagerImpl) RecreatePodContainer(pod *apiObject.Pod) error {
	return nil
}

func (*podManagerImpl) ExecPodContainer(pod *apiObject.Pod) error {
	return nil
}
