package netRequest

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	ContentType = "Content-Type"
	JSONType    = ""
)

// PostRequestByTarget 发送Post请求的工具函数，其中target可以是任意的数据类
func PostRequestByTarget(url string, target interface{}) (int, interface{}, error) {
	// 将目标对象序列化
	jsonData, err := json.Marshal(target)
	if err != nil {
		return 0, nil, err
	}

	response, err := http.Post(url, ContentType, bytes.NewBuffer(jsonData))
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
