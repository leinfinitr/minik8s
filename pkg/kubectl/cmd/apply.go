package cmd

import (
	"fmt"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/pkg/kubectl/translator"
	httprequest "minik8s/tools/httpRequest"
	"minik8s/tools/log"
	"net/http"
	"os"
	"strings"

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
	Hpa        ApplyObject = "Hpa"
)

func applyHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.ErrorLog("You must specify the type of resource to apply.")
		os.Exit(1)
	}
	if len(args) > 1 {
		log.ErrorLog("You may only specify one resource type to apply.")
		os.Exit(1)
	}
	if len(args) == 1 {
		resourceType := args[0]
		fileInfo, err := os.Stat(resourceType)
		if err != nil {
			log.ErrorLog("The resource type specified does not exist.")
			os.Exit(1)
		}
		if fileInfo.IsDir() {
			log.ErrorLog("The resource type specified is a directory.")
			os.Exit(1)
		}
		content, err := os.ReadFile(resourceType)
		if err != nil {
			log.ErrorLog("Could not read the file specified.")
			os.Exit(1)
		}
		kind, err := translator.FetchApiObjFromYaml(content)
		if err != nil {
			log.ErrorLog("Could not fetch the kind from the yaml file.")
			os.Exit(1)
		}
		switch kind {
		case "Pod":
			PodHandler(content)
		case "Service":
			ServiceHandler(content)
		case "Deployment":
			DeploymentHandler(content)
		case "Hpa":
			HpaHandler(content)
		case "ReplicaSet":
			ReplicaSetHandler(content)
		default:
			log.ErrorLog("The kind specified is not supported.")
			os.Exit(1)
		}
	}
}
func PodHandler(content []byte) {
	var pod apiObject.Pod
	err := translator.ParseApiObjFromYaml(content, &pod)
	if err != nil {
		log.ErrorLog("Could not unmarshal the yaml file.")
		os.Exit(1)
	}
	if pod.Metadata.Namespace == "" {
		pod.Metadata.Namespace = "default"
	}
	if pod.Metadata.Name == "" {
		log.ErrorLog("The name of the pod is required.")
		os.Exit(1)
	}
	url := config.APIServerURL() + config.PodsURI
	url = strings.Replace(url, config.NameSpaceReplace, pod.Metadata.Namespace, -1)
	log.DebugLog("Post " + url)
	resp, err := httprequest.PostObjMsg(url, pod)
	if err != nil {
		log.ErrorLog("Could not post the object message." + err.Error())
		os.Exit(1)
	}
	ApplyResultDisplay(Pod, resp)
}

func ServiceHandler(content []byte) {
	var service apiObject.Service
	err := translator.ParseApiObjFromYaml(content, &service)
	if err != nil {
		log.ErrorLog("Could not unmarshal the yaml file.")
		os.Exit(1)
	}
	if service.Metadata.Namespace == "" {
		service.Metadata.Namespace = "default"
	}
	if service.Metadata.Name == "" {
		log.ErrorLog("The name of the service is required.")
		os.Exit(1)
	}
	url := config.APIServerURL() + config.ServicesURI
	url = strings.Replace(url, config.NameSpaceReplace, service.Metadata.Namespace, -1)
	log.DebugLog("PUT " + url)
	resp, err := httprequest.PostObjMsg(url, service)
	if err != nil {
		log.ErrorLog("Could not post the object message." + err.Error())
		os.Exit(1)
	}
	ApplyResultDisplay(Service, resp)
}

func DeploymentHandler(content []byte) {

}

func HpaHandler(content []byte) {
	var hpa apiObject.HPA
	err := translator.ParseApiObjFromYaml(content, &hpa)
	if err != nil {
		log.ErrorLog("Could not unmarshal the yaml file.")
		os.Exit(1)
	}
	if hpa.Metadata.Namespace == "" {
		hpa.Metadata.Namespace = "default"
	}
	if hpa.Metadata.Name == "" {
		log.ErrorLog("The name of the hpa is required.")
		os.Exit(1)
	}
	url := config.APIServerURL() + config.HpasURI
	url = strings.Replace(url, config.NameSpaceReplace, hpa.Metadata.Namespace, -1)
	log.DebugLog("PUT " + url)
	resp, err := httprequest.PostObjMsg(url, hpa)
	if err != nil {
		log.ErrorLog("Could not post the object message." + err.Error())
		os.Exit(1)
	}
	ApplyResultDisplay(Hpa, resp)
}

func ReplicaSetHandler(content []byte) {
	var replicaSet apiObject.ReplicaSet
	err := translator.ParseApiObjFromYaml(content, &replicaSet)
	if err != nil {
		log.ErrorLog("Could not unmarshal the yaml file.")
		os.Exit(1)
	}
	if replicaSet.Metadata.Namespace == "" {
		replicaSet.Metadata.Namespace = "default"
	}
	if replicaSet.Metadata.Name == "" {
		log.ErrorLog("The name of the replicaSet is required.")
		os.Exit(1)
	}
	url := config.APIServerURL() + config.ReplicaSetURI
	url = strings.Replace(url, config.NameSpaceReplace, replicaSet.Metadata.Namespace, -1)
	log.DebugLog("PUT " + url)
	resp, err := httprequest.PostObjMsg(url, replicaSet)
	if err != nil {
		log.ErrorLog("Could not post the object message." + err.Error())
		os.Exit(1)
	}
	ApplyResultDisplay(ReplicaSet, resp)
}
func ApplyResultDisplay(kind ApplyObject, resp *http.Response) {
	if resp.StatusCode == http.StatusCreated {
		fmt.Printf("%s created\n", kind)
	} else if resp.StatusCode == http.StatusOK {
		fmt.Printf("%s updated\n", kind)
	} else {
		fmt.Printf("%s failed\n", kind)
	}
}
