package httprequest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"errors"
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

func DelObjMsg(url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetObjMsg(url string,obj interface{},kind string)(*http.Response,error){	
	res,err := http.Get(url)
	if err != nil {
		return nil,err
	}
	var result map[string]interface{}
	decoder := json.NewDecoder(res.Body)
	err2 := decoder.Decode(&result)
	if err2 != nil {
		return nil,err2
	}
	data,ok := result[kind]
	if !ok {
		return nil,errors.New("no such key")
	}
	err3:= json.Unmarshal([]byte(data.(string)),obj)
	if err3 != nil {
		return nil,err3
	}
	return res,nil
}