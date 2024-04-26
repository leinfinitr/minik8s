package netrequest

import (
	"bytes"
	"encoding/json"
	"minik8s/pkg/k8stype"
	"net/http"
)

// Put请求
func PutRequestByTarget(uri string, target interface{}) (int, interface{}, error) {
	jsonData, err := json.Marshal(target)
	if err != nil {
		return 0, nil, err
	}

	request, err := http.NewRequest("PUT", uri, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, nil, err
	}
	request.Header.Set("Content-Type", k8stype.ContentType)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return 0, nil, err
	}
	defer response.Body.Close()

	var bodyJson interface{}
	if err := json.NewDecoder(response.Body).Decode(&bodyJson); err != nil {
		return 0, nil, err
	}

	return response.StatusCode, bodyJson, nil
}
