package cmd

import (
	"fmt"
	httprequest "minik8s/internal/pkg/httpRequest"
	reflectprop "minik8s/internal/pkg/reflectProp"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/pkg/kubectl/translator"
	"os"
	"reflect"
	"net/http"
	"github.com/spf13/cobra"
)

var deletedCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource by filename or stdin",
	Long:  "Delete resources by filenames, stdin, resources and names, or by resources and label selector",
	Run:   deleteHandler,
}

func deleteHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: You must specify the type of resource to delete.")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("Error: You may only specify one resource type to delete.")
		os.Exit(1)
	}
	if len(args) == 1 {
		resourceType := args[0]
		fileInfo, err := os.Stat(resourceType)
		if err != nil {
			fmt.Println("Error: The resource type specified does not exist.")
			os.Exit(1)
		}
		if fileInfo.IsDir() {
			fmt.Println("Error: The resource type specified is a directory.")
			os.Exit(1)
		}
		content, err := os.ReadFile(resourceType)
		if err != nil {
			fmt.Println("Error: Could not read the file specified.")
			os.Exit(1)
		}
		kind, err := translator.FetchApiObjFromYaml(content)
		if err != nil {
			fmt.Println("Error: Could not fetch the kind from the yaml file.")
			os.Exit(1)
		}
		var ActualType reflect.Type
		switch kind {
		case "Pod":
			ActualType = reflect.TypeOf(apiObject.Pod{})
		case "Service":
			ActualType = reflect.TypeOf(apiObject.Service{})
		case "ReplicaSet":
			ActualType = reflect.TypeOf(apiObject.ReplicaSet{})
		default:
			ActualType = nil
		}
		if ActualType == nil {
			fmt.Println("Error: The resource type specified is not supported.")
			os.Exit(1)
		}
		obj := reflect.New(ActualType).Interface()
		err = translator.ParseApiObjFromYaml(content, obj)
		if err != nil {
			fmt.Println("Error: Could not parse the yaml file.")
			os.Exit(1)
		}
		namespace := reflectprop.GetObjNamespace(obj)
		name := reflectprop.GetObjName(obj)
		if name == "" {
			fmt.Println("Error: Could not get the name of the resource.")
			os.Exit(1)
		}
		url := config.APIServerUrl() + "/api/v1/namespaces/" + namespace + "/" + kind + "/" + name
		resp, err := httprequest.DelObjMsg(url)
		if err != nil {
			fmt.Println("Error: Could not delete the object.")
			os.Exit(1)
		}
		
		DeleteResultDisplay(kind, resp)
	}
}

func DeleteResultDisplay(kind string, resp *http.Response) {
	if(resp.StatusCode == 200) {
		fmt.Println(kind + " deleted successfully.")
	} else {
		fmt.Println("Error: Could not delete the " + kind + ".")
	}
}
