package specctlrs

import (
	"errors"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/executor"
	httprequest "minik8s/tools/httpRequest"
	"minik8s/tools/log"
	netRequest "minik8s/tools/netRequest"
	stringops "minik8s/tools/stringops"
	"net/http"
	"strings"
	"time"
)

type ReplicaSetController interface {
	Run()
}
type ReplicaSetControllerImpl struct {
}

var (
	ReplicaControllerDelay   = 3 * time.Second
	ReplicaControllerTimeGap = []time.Duration{10 * time.Second}
)

func NewReplicaController() (ReplicaSetController, error) {
	return &ReplicaSetControllerImpl{}, nil
}

func (rc *ReplicaSetControllerImpl) Run() {
	// 定期执行
	executor.ExecuteInPeriod(ReplicaControllerDelay, ReplicaControllerTimeGap, rc.syncReplicaSet)
}


func GetAllReplicaSetsFromAPIServer() ( replicaSets []apiObject.ReplicaSet,err error) {
	url := config.APIServerURL() + config.ReplicaSetsURI
	res, err := httprequest.GetObjMsg(url, &replicaSets, "data")
	if err != nil {
		log.ErrorLog("GetAllReplicaSetsFromAPIServer: " + err.Error())
		return replicaSets, err 
	}
	if res.StatusCode != 200 {
		log.ErrorLog("GetAllReplicaSetsFromAPIServer: " + res.Status)
		return replicaSets,err
	}
	return replicaSets,nil
}
func (rc *ReplicaSetControllerImpl) syncReplicaSet() {
	var pods []apiObject.Pod
	// 1. 获取所有的Pod
	pods,err := GetAllPodsFromAPIServer()
	if err != nil {
		log.ErrorLog("syncReplicaSet: " + err.Error())
		return
	}
	// 2. 获取所有的ReplicaSet
	var replicaSets []apiObject.ReplicaSet
	replicaSets,err = GetAllReplicaSetsFromAPIServer()
	if err != nil {
		log.ErrorLog("syncReplicaSet: " + err.Error())
		return
	}
	replicaMapping := make(map[string]string, 0)
	for _, rs := range replicaSets {
		key := rs.Metadata.Namespace + "/" + rs.Metadata.Name
		replicaMapping[key] = rs.Metadata.UUID
	}

	for _, rs := range replicaSets {
		selectedPods := []apiObject.Pod{}
		for _, pod := range pods {
			if PodsMatched(pod, rs.Spec.Selector) {
				selectedPods = append(selectedPods, pod)
			}
		}
		if len(selectedPods) < int(rs.Spec.Replicas) {
			log.InfoLog("syncReplicaSet: " + rs.Metadata.Name + " need to scale")
			// 3. 如果Pod数量不足，则创建Pod
			err := rc.IncreaseReplicas(&rs.Metadata, &rs.Spec.Template, int(rs.Spec.Replicas)-len(selectedPods))
			if err != nil {
				log.ErrorLog("syncReplicaSet: " + err.Error())
			}
		} else if len(selectedPods) > int(rs.Spec.Replicas) {
			log.InfoLog("syncReplicaSet: " + rs.Metadata.Name + " need to scale")
			// 4. 如果Pod数量过多，则删除Pod
			err := rc.DecreaseReplicas(selectedPods, len(selectedPods)-int(rs.Spec.Replicas))
			if err != nil {
				log.ErrorLog("syncReplicaSet: " + err.Error())
			}
		}
		rc.UpdateStatus(&rs, selectedPods)
	}
	//对于已经删除的replicaSet，需要删除对应的pod
	for _, pod := range pods {
		if pod.Metadata.Labels[apiObject.PodReplicaUUID] != "" {
			if pod.Metadata.Labels[apiObject.PodReplicaNamespace] == "" || pod.Metadata.Labels[apiObject.PodReplicaName] == "" {
				continue
			}
			key := pod.Metadata.Labels[apiObject.PodReplicaNamespace] + "/" + pod.Metadata.Labels[apiObject.PodReplicaName]
			if _, ok := replicaMapping[key]; !ok {
				rc.DecreaseReplicas([]apiObject.Pod{pod}, 1)
			}
		}
	}

}


func (rc *ReplicaSetControllerImpl) IncreaseReplicas(replicaMeta *apiObject.ObjectMeta, pod *apiObject.PodTemplateSpec, num int) error {
	new_pod := apiObject.Pod{}
	new_pod.Metadata = pod.Metadata
	new_pod.Kind = apiObject.PodType
	new_pod.APIVersion = "v1"
	new_pod.Spec = pod.Spec
	new_pod.Metadata.Labels[apiObject.PodReplicaName] = replicaMeta.Name
	new_pod.Metadata.Labels[apiObject.PodReplicaNamespace] = replicaMeta.Namespace
	new_pod.Metadata.Labels[apiObject.PodReplicaUUID] = replicaMeta.UUID

	originalPodName := new_pod.Metadata.Name

	originalContainerNames := make([]string, 0)

	for _, container := range new_pod.Spec.Containers {
		originalContainerNames = append(originalContainerNames, container.Name)
	}

	url := config.APIServerURL() + config.PodsURI

	errStr := ""
	for i := 0; i < num; i++ {
		new_pod.Metadata.Name = originalPodName + "-" + stringops.GenerateRandomString(5)

		// 修改container的name
		for idx := range new_pod.Spec.Containers {
			new_pod.Spec.Containers[idx].Name = originalContainerNames[idx] + "-" + stringops.GenerateRandomString(5)
		}

		url = strings.Replace(url, config.NameSpaceReplace, pod.Metadata.Namespace, -1)
		code, _, err := netRequest.PostRequestByTarget(url, &new_pod)

		if err != nil {
			log.ErrorLog("replicaController: " + "AddPodsNums error: " + err.Error())
			errStr += err.Error()
		}

		if code != http.StatusCreated {
			log.ErrorLog("replicaController: " + "AddPodsNums code is not 201")
			errStr += "code is not 200"
		}
	}

	if errStr != "" {
		return errors.New(errStr)
	}

	return nil
}

func (rc *ReplicaSetControllerImpl) DecreaseReplicas(matchedPods []apiObject.Pod, num int) error {
	if len(matchedPods) < num {
		return errors.New("matchedPods is less than num")
	}

	for i := 0; i < num; i++ {
		pod := matchedPods[i]
		url := config.APIServerURL() + config.PodURI
		url = strings.Replace(url, config.NameSpaceReplace, pod.Metadata.Namespace, -1)
		url = strings.Replace(url, config.NameReplace, pod.Metadata.Name, -1)
		code, err := netRequest.DelRequest(url)

		if err != nil {
			log.ErrorLog("replicaController: " + "DeletePodsNums error: " + err.Error())
		}

		if code != http.StatusOK {
			log.ErrorLog("replicaController: " + "DeletePodsNums code is not 200")
		}
	}

	return nil
}

func (rc *ReplicaSetControllerImpl) UpdateStatus(replicaSet *apiObject.ReplicaSet, pods []apiObject.Pod) error {
	newReplicaStatus := apiObject.ReplicaSetStatus{}
	newReplicaStatus.Conditions = []apiObject.ReplicaSetCondition{}
	numsReady := 0
	for _, pod := range pods {
		if pod.Status.Phase == apiObject.PodRunning {
			numsReady++
		}
		newReplicaStatus.Conditions = append(newReplicaStatus.Conditions, apiObject.ReplicaSetCondition{
			Type:               apiObject.PodType,
			Status:             pod.Status.Phase,
			LastTransitionTime: time.Now().String(),
		})
	}
	newReplicaStatus.Replicas = replicaSet.Spec.Replicas
	newReplicaStatus.ReadyReplicas = int32(numsReady)

	url := config.APIServerURL() + config.ReplicaSetStatusURI
	url = strings.Replace(url, config.NameSpaceReplace, replicaSet.Metadata.Namespace, -1)
	url = strings.Replace(url, config.NameReplace, replicaSet.Metadata.Name, -1)
	code, _, err := netRequest.PutRequestByTarget(url, &newReplicaStatus)
	if err != nil {
		log.ErrorLog("replicaController: " + "UpdateStatus error: " + err.Error())
		return err
	}
	if code != http.StatusOK {
		log.ErrorLog("replicaController: " + "UpdateStatus code is not 200")
		return errors.New("UpdateStatus code is not 200")
	}
	return nil
}