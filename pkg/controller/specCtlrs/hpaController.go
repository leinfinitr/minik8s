package specctlrs

import (
	"encoding/json"
	"math"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/executor"
	"minik8s/tools/log"
	netRequest "minik8s/tools/netRequest"
	stringops "minik8s/tools/stringops"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type HpaController interface {
	Run()
}

type HpaControllerImpl struct {
}

var (
	HpaControllerDelay   = 0 * time.Second
	HpaControllerTimeGap = []time.Duration{10 * time.Second}
)

func NewHpaController() (HpaController, error) {
	return &HpaControllerImpl{}, nil
}

func (hc *HpaControllerImpl) Run() {
	// 定期执行
	executor.ExecuteInPeriod(HpaControllerDelay, HpaControllerTimeGap, hc.syncHpa)
}

func GetAllHpasFromAPIServer() (hpas []apiObject.HPA, err error) {
	url := config.APIServerURL() + config.GlobalHpaURI
	res, err := http.Get(url)
	if err != nil {
		log.ErrorLog("GetAllHpasFromAPIServer: " + err.Error())
		return hpas, err
	}
	err = json.NewDecoder(res.Body).Decode(&hpas)
	if err != nil {
		log.ErrorLog("GetAllHpasFromAPIServer: " + err.Error())
		return hpas, err
	}
	return hpas, nil
}

func (hc *HpaControllerImpl) syncHpa() {
	hpas, err := GetAllHpasFromAPIServer()
	if err != nil {
		log.ErrorLog("syncHpa: " + err.Error())
		return
	}
	pods, err := GetAllPodsFromAPIServer()
	if err != nil {
		log.ErrorLog("syncHpa: " + err.Error())
		return
	}
	hpaMapping := make(map[string]string, 0)
	for _, hpa := range hpas {
		key := hpa.Metadata.Namespace + "/" + hpa.Metadata.Name
		hpaMapping[key] = key
	}

	for _, rs := range hpas {
		go hc.handleHPA(rs, pods)
	}
}

func (hc *HpaControllerImpl) handleHPA(hpa apiObject.HPA, pods []apiObject.Pod) {
	selectedPods := make([]apiObject.Pod, 0)
	for _, pod := range pods {
		if PodsMatched(pod, hpa.Spec.Selector) {
			selectedPods = append(selectedPods, pod)
		}
	}
	log.DebugLog("selectedPods: " + strconv.Itoa(len(selectedPods)))
	HpaControllerTimeGap = []time.Duration{min(HpaControllerTimeGap[0], time.Duration(hpa.Spec.AdjustInterval))}
	hpa.Status.CurrentReplicas = int32(len(selectedPods))
	if hpa.Status.CurrentReplicas == 0 {
		log.ErrorLog("handleHPA: " + hpa.Metadata.Namespace + "/" + hpa.Metadata.Name + " no pod selected")
		return
	}
	if hpa.Status.CurrentReplicas < int32(hpa.Spec.MinReplicas) {
		err := hc.AddOnePod(hpa, selectedPods[0])
		if err != nil {
			log.ErrorLog("handleHPA: " + hpa.Metadata.Namespace + "/" + hpa.Metadata.Name + " add one pod failed")
		}
		return
	}
	if hpa.Status.CurrentReplicas > int32(hpa.Spec.MaxReplicas) {
		err := hc.DeleteOnePod(selectedPods[0])
		if err != nil {
			log.ErrorLog("handleHPA: " + hpa.Metadata.Namespace + "/" + hpa.Metadata.Name + " delete one pod failed")
		}
		return
	}
	avgCPU := hc.CalAvgCpuCost(selectedPods)
	avgMem := hc.CalAvgMemCost(selectedPods)

	expectedNm := hc.CalExpectedReplicas(avgCPU, avgMem, hpa)
	if int32(expectedNm) > hpa.Status.CurrentReplicas {
		err := hc.AddOnePod(hpa, selectedPods[0])
		if err != nil {
			log.ErrorLog("handleHPA: " + hpa.Metadata.Namespace + "/" + hpa.Metadata.Name + " add one pod failed")
		}
	}
	log.DebugLog("current replicas: " + strconv.Itoa(int(hpa.Status.CurrentReplicas)) + " expected replicas: " + strconv.Itoa(expectedNm))
	if int32(expectedNm) < hpa.Status.CurrentReplicas {
		err := hc.DeleteOnePod(selectedPods[0])
		if err != nil {
			log.ErrorLog("handleHPA: " + hpa.Metadata.Namespace + "/" + hpa.Metadata.Name + " delete one pod failed")
		}
	}

	//下一次Update时候会将数量更新正确
	hpa.Status.CurCPUPercent = avgCPU
	hpa.Status.CurMemoryPercent = avgMem

	err := hc.UpdateStatus(hpa)
	if err != nil {
		log.ErrorLog("handleHPA: " + hpa.Metadata.Namespace + "/" + hpa.Metadata.Name + " update status failed")
	}
}

func (hc *HpaControllerImpl) AddOnePod(hpa apiObject.HPA, pod apiObject.Pod) error {
	log.InfoLog("AddOnePod: " + hpa.Metadata.Namespace + "/" + hpa.Metadata.Name + " add one pod")
	new_pod := pod
	new_pod.Metadata.Labels[apiObject.PodHpaName] = hpa.Metadata.Name
	new_pod.Metadata.Labels[apiObject.PodHpaNamespace] = hpa.Metadata.Namespace
	new_pod.Metadata.Labels[apiObject.PodHpaUUID] = hpa.Metadata.UUID
	for idx := range new_pod.Spec.Containers {
		new_pod.Spec.Containers[idx].Name = new_pod.Spec.Containers[idx].Name + "-" + stringops.GenerateRandomString(5)
	}
	new_pod.Metadata.Name = pod.Metadata.Name + "-" + stringops.GenerateRandomString(5)
	url := config.APIServerURL() + config.PodsURI
	url = strings.Replace(url, config.NameSpaceReplace, hpa.Metadata.Namespace, -1)
	code, _, err := netRequest.PostRequestByTarget(url, new_pod)
	if err != nil {
		log.ErrorLog("AddOnePod: " + hpa.Metadata.Namespace + "/" + hpa.Metadata.Name + " add one pod failed")
		return err
	}
	if code != 201 {
		log.ErrorLog("AddOnePod: " + hpa.Metadata.Namespace + "/" + hpa.Metadata.Name + " add one pod failed")
		return err
	}
	return nil
}

func (hc *HpaControllerImpl) DeleteOnePod(pod apiObject.Pod) error {
	url := config.APIServerURL() + config.PodURI
	url = strings.Replace(url, config.NameSpaceReplace, pod.Metadata.Namespace, -1)
	url = strings.Replace(url, config.NameReplace, pod.Metadata.Name, -1)
	code, err := netRequest.DelRequest(url)
	if err != nil {
		log.ErrorLog("DeleteOnePod: " + pod.Metadata.Namespace + "/" + pod.Metadata.Name + " delete one pod failed")
		return err
	}
	if code != 200 {
		log.ErrorLog("DeleteOnePod: " + pod.Metadata.Namespace + "/" + pod.Metadata.Name + " delete one pod failed")
		return err
	}
	return nil
}

func (hc *HpaControllerImpl) CalAvgCpuCost(pods []apiObject.Pod) float64 {
	totalCost := 0.0
	for _, pod := range pods {
		totalCost += pod.Status.CpuUsage
	}
	return totalCost / float64(len(pods))
}

func (hc *HpaControllerImpl) CalAvgMemCost(pods []apiObject.Pod) float64 {
	totalCost := 0.0
	for _, pod := range pods {
		totalCost += pod.Status.MemUsage
	}
	return totalCost / float64(len(pods))
}

func (hc *HpaControllerImpl) CalExpectedReplicas(avgCpuCost float64, avgMemCost float64, hpa apiObject.HPA) int {
	cpuUsedPer := avgCpuCost / hpa.Spec.Metrics.CPUPercent
	memUsedPer := avgMemCost / hpa.Spec.Metrics.MemoryPercent
	expectedNm := int(math.Max(cpuUsedPer, memUsedPer) * float64(hpa.Status.CurrentReplicas))
	log.DebugLog("cpuUsedPer: " + strconv.FormatFloat(cpuUsedPer, 'f', 2, 64))
	log.DebugLog("memUsedPer: " + strconv.FormatFloat(memUsedPer, 'f', 2, 64))
	log.DebugLog("expectedNm: " + strconv.Itoa(expectedNm))

	if expectedNm < int(hpa.Spec.MinReplicas) {
		return int(hpa.Spec.MinReplicas)
	}
	if expectedNm > int(hpa.Spec.MaxReplicas) {
		return int(hpa.Spec.MaxReplicas)
	}
	return expectedNm
}

func (hc *HpaControllerImpl) UpdateStatus(hpa apiObject.HPA) error {
	url := config.APIServerURL() + config.HpaStatusURI
	url = strings.Replace(url, config.NameSpaceReplace, hpa.Metadata.Namespace, -1)
	url = strings.Replace(url, config.NameReplace, hpa.Metadata.Name, -1)
	code, _, err := netRequest.PutRequestByTarget(url, &hpa.Status)
	if err != nil {
		log.ErrorLog("UpdateStatus: " + hpa.Metadata.Namespace + "/" + hpa.Metadata.Name + " update status failed")
		return err
	}
	if code != 200 {
		log.ErrorLog("UpdateStatus: " + hpa.Metadata.Namespace + "/" + hpa.Metadata.Name + " update status failed")
		return err
	}
	return nil
}
