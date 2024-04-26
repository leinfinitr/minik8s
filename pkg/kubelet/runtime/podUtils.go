package runtime

import (
	"minik8s/pkg/apiObject"
)

type PodUtils interface {
	CreatePod(pod *apiObject.Pod) error
	StartPod(pod *apiObject.Pod) error
	RestartPod(pod *apiObject.Pod) error
	StopPod(pod *apiObject.Pod) error
	DeletePod(pod *apiObject.Pod) error
	RecreatePodContainer(pod *apiObject.Pod) error
	ExecPodContainer(pod *apiObject.Pod) error
}

func (r *RuntimeManager) CreatePod(pod *apiObject.Pod) error {
	/* Step-1: Build pause container*/

	/* Step-2: Build pod中所需的所有container*/
	return nil
}

func (r *RuntimeManager) StartPod(pod *apiObject.Pod) error {
	return nil
}

func (r *RuntimeManager) RestartPod(pod *apiObject.Pod) error {
	return nil
}

func (r *RuntimeManager) StopPod(pod *apiObject.Pod) error {
	return nil
}

func (r *RuntimeManager) DeletePod(pod *apiObject.Pod) error {
	return nil
}

func (r *RuntimeManager) RecreatePodContainer(pod *apiObject.Pod) error {
	return nil
}

func (r *RuntimeManager) ExecPodContainer(pod *apiObject.Pod) error {
	return nil
}
