package httprequest
import (
	"bytes"
	"net/http"
	"encoding/json"
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