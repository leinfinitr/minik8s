package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	httprequest "minik8s/internal/pkg/httpRequest"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe a resource by filename or stdin",
	Long:  "Describe resources by filenames, stdin, resources and names, or by resources and label selector",
	Run:   describeHandler,
}

func describeHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: You must specify the type of resource to describe.")
		os.Exit(1)
	}
	resourceType := args[0]
	if len(args)==1{
		namespace := "default"
		describeNamespace(namespace, resourceType)
	}else if len(args)==2{
		//describe [resource-type] [namespace]/[resource-name]
		namespace, resourceName := SplitNamespaceAndResourceName(args[1])
		if namespace == "" || resourceName == "" {
			fmt.Println("Error: The resource name must be in the format namespace/resource-name")
			os.Exit(1)
		}
		describeResource(namespace, resourceName, resourceType)
	}else{
		fmt.Println("Error: The resource name must be in the format namespace/resource-name")
		os.Exit(1)
	}
}

func describeNamespace(namespace string, resourceType string) {
	url := config.APIServerURL() + config.UriMapping[resourceType]
	url = strings.Replace(url, config.NameSpaceReplace, namespace, -1)
	var ActualType reflect.Type
	switch resourceType {
	case "Pod":
		ActualType = reflect.TypeOf(apiObject.Pod{})
	case "Service":
		ActualType = reflect.TypeOf(apiObject.Service{})
	case "ReplicaSet":
		ActualType = reflect.TypeOf(apiObject.ReplicaSet{})
	default:
		ActualType = nil
	}
	obj := reflect.New(ActualType).Interface()
	res, err := httprequest.GetObjMsg(url, obj, "data")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	if res.StatusCode != 200 {
		fmt.Println("Error: Failed to get nodes")
		os.Exit(1)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading body: ", err)
		os.Exit(1)
	}

	// 使用 json.Indent 对 JSON 数据进行格式化
	var out bytes.Buffer
	err = json.Indent(&out, bodyBytes, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON: ", err)
		os.Exit(1)
	}

	fmt.Println("Resource Details:")
	fmt.Println(out.String())
}

func describeResource(namespace string, resourceName string, resourceType string) {
	url := config.APIServerURL() + config.UriSpecMapping[resourceType]
	url = strings.Replace(url, config.NameSpaceReplace, namespace, -1)
	url = strings.Replace(url, config.NameReplace, resourceName, -1)
	var ActualType reflect.Type
	switch resourceType {
	case "Pod":
		ActualType = reflect.TypeOf(apiObject.Pod{})
	case "Service":
		ActualType = reflect.TypeOf(apiObject.Service{})
	case "ReplicaSet":
		ActualType = reflect.TypeOf(apiObject.ReplicaSet{})
	default:
		ActualType = nil
	}
	obj := reflect.New(ActualType).Interface()
	res, err := httprequest.GetObjMsg(url, obj, "data")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	if res.StatusCode != 200 {
		fmt.Println("Error: Failed to get nodes")
		os.Exit(1)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading body: ", err)
		os.Exit(1)
	}

	// 使用 json.Indent 对 JSON 数据进行格式化
	var out bytes.Buffer
	err = json.Indent(&out, bodyBytes, "", "  ")
	if err != nil {
		fmt.Println("Error formatting JSON: ", err)
		os.Exit(1)
	}

	fmt.Println("Resource Details:")
	fmt.Println(out.String())
}

func SplitNamespaceAndResourceName(resource string) (string, string) {
	split := strings.Split(resource, "/")
	if len(split) != 2 {
		fmt.Println("Error: The resource name must be in the format namespace/resource-name")
		os.Exit(1)
	}
	return split[0], split[1]
}