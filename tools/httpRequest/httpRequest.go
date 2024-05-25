package httprequest

import (
	"bytes"
	"encoding/json"
	"errors"

	// "fmt"
	"net/http"
)

func PostObjMsg(url string, obj interface{}) (*http.Response, error) {
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func DelMsg(url string, obj interface{}) (*http.Response, error) {
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func PutObjMsg(url string, obj interface{}) (*http.Response, error) {
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetObjMsg(url string, obj interface{}, kind string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	decoder := json.NewDecoder(res.Body)
	err2 := decoder.Decode(&result)
	if err2 != nil {
		return nil, err2
	}
	data, ok := result[kind]
	if !ok {
		return nil, errors.New("no such key")
	}
	err3 := json.Unmarshal([]byte(data.(string)), obj)
	if err3 != nil {
		return nil, err3
	}
	return res, nil
}

func GetMsg(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return res, nil
}
