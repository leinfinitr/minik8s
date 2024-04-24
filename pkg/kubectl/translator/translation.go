package translator

import (
	"gopkg.in/yaml.v3"
	"fmt"
)

func FetchApiObjFromYaml(yamlBytes []byte) (string, error) {
	var res map[string]interface{}
	err := yaml.Unmarshal(yamlBytes, &res)
	if err != nil {
		return "", err
	}
	if res["kind"] == nil {
		return "", fmt.Errorf("kind is required")
	}
	return res["kind"].(string), nil
}

func ParseApiObjFromYaml(yamlBytes []byte, obj interface{}) error {
	err := yaml.Unmarshal(yamlBytes, obj)
	if err != nil {
		return err
	}
	return nil
}