package netrequest

import (
	"bytes"
	"encoding/json"
	"minik8s/pkg/k8stype"
	"net/http"
)

/* 发送Post请求的工具函数，其中target可以是任意的数据类*/
func PostRequestByTarget(uri string, target interface{}) (int, interface{}, error) {
	// 将目标对象序列化
	jsonData, err := json.Marshal(target)
	if err != nil {
		return 0, nil, err
	}

	response, err := http.Post(uri, k8stype.ContentType, bytes.NewBuffer(jsonData))
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

func PostString(uri string, str string) (*http.Response, error) {
	cli := http.Client{}
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader([]byte(str)))
	if err != nil {
		return nil, err
	}
	return cli.Do(req)
}
