package cmd

import (
	"fmt"
	"minik8s/internal/pkg/httpRequest"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/pkg/kubectl/translator"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a configuration to a resource by filename or stdin",
	Long:  "Apply a configuration to a resource by filename or stdin",
	Run:   applyHandler,
}

type ApplyObject string

const (
	Pod        ApplyObject = "Pod"
	Service    ApplyObject = "Service"
	Deployment ApplyObject = "Deployment"
	ReplicaSet ApplyObject = "ReplicaSet"
)

func applyHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: You must specify the type of resource to apply.")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("Error: You may only specify one resource type to apply.")
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
		switch kind {
		case "Pod":
			PodHandler(content)
		case "Service":
			ServiceHandler(content)
		case "Deployment":
			DeploymentHandler(content)
		default:
			fmt.Println("Error: The kind specified is not supported.")
			os.Exit(1)
		}
	}
}
func PodHandler(content []byte) {
	var pod apiObject.Pod
	err := translator.ParseApiObjFromYaml(content, &pod)
	if err != nil {
		fmt.Println("Error: Could not unmarshal the yaml file.")
		os.Exit(1)
	}
	url := config.APIServerUrl() + "/api/v1/namespaces/" + pod.Metadata.Namespace + "/pods"
	resp, err := httprequest.PostObjMsg(url, pod)
	if err != nil {
		fmt.Println("Error: Could not post the object message.")
		os.Exit(1)
	}
	ResultDisplay(Pod, resp)
}

func ServiceHandler(content []byte) {
	// var service apiObject.Service
	// err := translator.ParseApiObjFromYaml(content, &service)
	// if err != nil {
	// 	fmt.Println("Error: Could not unmarshal the yaml file.")
	// 	os.Exit(1)
	// }

}

func DeploymentHandler(content []byte) {

}

func ResultDisplay(kind ApplyObject, resp *http.Response) {
	if resp.StatusCode == http.StatusCreated {
		fmt.Printf("%s created\n", kind)
	} else if resp.StatusCode == http.StatusOK {
		fmt.Printf("%s updated\n", kind)
	} else {
		fmt.Printf("%s failed\n", kind)
	}
}
