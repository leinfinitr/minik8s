package specctlrs

import (
	"encoding/json"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"
	"net/http"
)
func GetAllPodsFromAPIServer() (pods []apiObject.Pod,err error) {
	url := config.APIServerURL() + config.PodsGlobalURI
	res, err := http.Get(url)
	if err != nil {
		log.ErrorLog("GetAllPodsFromAPIServer: " + err.Error())
		return pods,err
	}
	if res.StatusCode != 200 {
		log.ErrorLog("GetAllPodsFromAPIServer: " + res.Status)
		return pods,err
	}
	err = json.NewDecoder(res.Body).Decode(&pods)
	if err != nil {
		log.ErrorLog("GetAllPodsFromAPIServer: " + err.Error())
		return pods,err
	}
	return pods,nil
}

func PodsMatched(pod apiObject.Pod, selector map[string]string) bool {
	for k, v := range selector {
		if pod.Metadata.Labels[k] != v {
			return false
		}
	}
	return true
}
