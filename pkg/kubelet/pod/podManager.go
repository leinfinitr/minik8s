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
	log.DebugLog("[PodManager] Arrived into AddPod")
	uuid := pod.GetPodUUID()
	if _, ok := p.PodMapByUUID[uuid]; ok {
		log.ErrorLog("Pod has been built already")
		return errors.New("pod message has been handled")
	}

	p.PodMapByUUID[uuid] = pod
	pod.Status.Phase = apiObject.Pod_Building

	go func() {
		err := p.AddPodHandler(pod)
		if err != nil {
			log.ErrorLog("AddPodHandler error: " + err.Error())
			pod.Status.Phase = apiObject.Pod_Unknown
		} else {
			log.InfoLog("AddPodHandler success")
			pod.Status.Phase = apiObject.Pod_Succeeded
		}
	}()

	return nil
}

func (p *podManagerImpl) DeletePod(pod *apiObject.Pod) error {
	log.DebugLog("[PodManager] Arrived into DeletePod")
	uuid := pod.GetPodUUID()
	if _, ok := p.PodMapByUUID[uuid]; !ok {
		log.ErrorLog("Pod has been deleted already")
		return errors.New("pod message has been handled")
	}

	go func() {
		err := p.DeletePodHandler(pod)
		if err != nil {
			log.ErrorLog("DeletePodHandler error: " + err.Error())
		} else {
			log.InfoLog("DeletePodHandler success")
		}
		delete(p.PodMapByUUID, uuid)
	}()
	return nil
}

func (p *podManagerImpl) StartPod(pod *apiObject.Pod) error {
	var msg string
	log.DebugLog("[PodManager] Arrived into StartPod")
	uuid := pod.GetPodUUID()
	if _, ok := p.PodMapByUUID[uuid]; !ok {
		msg = "Pod cann't be found"
		log.ErrorLog(msg)
		return errors.New(msg)
	}

	// 需要对Pod的不同状况进行处理
	switch pod.Status.Phase {
	case apiObject.Pod_Succeeded:
		go func() {
			err := p.StartPodHandler(pod)
			if err != nil {
				log.ErrorLog("StartPodHandler error: " + err.Error())
			} else {
				log.InfoLog("StartPodHandler success")
			}
			pod.Status.Phase = apiObject.Pod_Running
		}()
		return nil
	case apiObject.Pod_Running:
		log.DebugLog("Pod has been running")
		return nil
	case apiObject.Pod_Building:
		msg = "Pod has not been built now! "
	default:
		msg = "Pod is not ready to start "
	}

	log.ErrorLog(msg)
	return errors.New(msg)
}

func (p *podManagerImpl) StopPod(pod *apiObject.Pod) error {
	log.DebugLog("[PodManager] Arrived into StopPod")
	uuid := pod.GetPodUUID()
	if _, ok := p.PodMapByUUID[uuid]; !ok {
		msg := "Pod cann't be found"
		log.ErrorLog(msg)
		return errors.New(msg)
	} else if pod.Status.Phase == apiObject.Pod_Running {
		go func() {
			err := p.StopPodHandler(pod)
			if err != nil {
				log.ErrorLog("StopPodHandler error: " + err.Error())
				pod.Status.Phase = apiObject.Pod_Unknown
			} else {
				log.InfoLog("StopPodHandler success")
				// 回退到所有容器都被创建好的状态
				pod.Status.Phase = apiObject.Pod_Succeeded
			}
		}()
		return nil
	} else {
		msg := "Status Error!"
		log.ErrorLog(msg)
		return errors.New(msg)
	}

}

func (p *podManagerImpl) RestartPod(pod *apiObject.Pod) error {
	log.DebugLog("[PodManager] Arrived into RestartPod")
	uuid := pod.GetPodUUID()
	if _, ok := p.PodMapByUUID[uuid]; !ok {
		msg := "Pod cann't be found"
		log.ErrorLog(msg)
		return errors.New(msg)
	} else if pod.Status.Phase == apiObject.Pod_Succeeded || pod.Status.Phase == apiObject.Pod_Failed ||
		pod.Status.Phase == apiObject.Pod_Running {
		go func() {
			err := p.StopPodHandler(pod)
			if err != nil {
				log.ErrorLog("RestartPodHandler error: " + err.Error())
				pod.Status.Phase = apiObject.Pod_Unknown
			} else {
				log.InfoLog("RestartPodHandler success")
				pod.Status.Phase = apiObject.Pod_Running
			}
			// 回退到所有容器都被创建好的状态
			pod.Status.Phase = apiObject.Pod_Succeeded
		}()
		return nil
	} else {
		msg := "Status Error!"
		log.ErrorLog(msg)
		return errors.New(msg)
	}

}

func (p *podManagerImpl) DeletePodByUUID(pod *apiObject.Pod) error {
	return nil
}

func (p *podManagerImpl) RecreatePodContainer(pod *apiObject.Pod) error {
	return nil
}

func (p *podManagerImpl) ExecPodContainer(pod *apiObject.Pod) error {
	return nil
}
