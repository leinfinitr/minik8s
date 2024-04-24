package runtime

import (
	"minik8s/pkg/apiObject"
)

type PodUtils interface {
	createPod(pod *apiObject.Pod) error
	startPod(pod *apiObject.Pod) error
	restartPod(pod *apiObject.Pod) error
	stopPod(pod *apiObject.Pod) error
	deletePod(pod *apiObject.Pod) error
	recreatePodContainer(pod *apiObject.Pod) error
	execPodContainer(pod *apiObject.Pod) error
}

func (r *RuntimeManager) CreatePod(pod *apiObject.Pod) error {
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
