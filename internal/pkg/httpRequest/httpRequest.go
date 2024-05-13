package httprequest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	// "fmt"
	"net/http"
)

func PostObjMsg(url string, obj interface{}) (*http.Response, error) {
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	// fmt.Println(string(jsonStr))
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

func GetObjMsg(url string,obj interface{},kind string)(*http.Response,error){	
	res,err := http.Get(url)
	if err != nil {
		return nil,err
	}
	var result map[string]interface{}
	fmt.Println(res.Body)
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

func GetObjMsg2(url string, obj interface{}, kind string) (*http.Response, error) {
    res, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    var result map[string]interface{}
    decoder := json.NewDecoder(res.Body)
    err = decoder.Decode(&result)
    if err != nil {
        return nil, err
    }

    data, ok := result[kind]
    if !ok {
        return nil, errors.New("no such key")
    }

    switch v := data.(type) {
    case []interface{}:
        objValue := reflect.ValueOf(obj)
        if objValue.Kind() != reflect.Ptr || objValue.Elem().Kind() != reflect.Slice {
            return nil, errors.New("obj should be a pointer to a slice")
        }
        elementType := objValue.Elem().Type().Elem()
        nodes := reflect.MakeSlice(reflect.SliceOf(elementType), 0, 0)

        for _, item := range v {
            itemStr, ok := item.(string)
            if !ok {
                return nil, errors.New("invalid data type")
            }
            newItem := reflect.New(elementType).Interface()
            err = json.Unmarshal([]byte(itemStr), newItem)
            if err != nil {
                return nil, err
            }
            nodes = reflect.Append(nodes, reflect.ValueOf(newItem).Elem())
        }

        objValue.Elem().Set(nodes)
    default:
        return nil, errors.New("unexpected data type")
    }

    return res, nil
}
