package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"minik8s/pkg/config"

	httprequest "minik8s/tools/httpRequest"
)

var deletedCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource by namespace, type and name",
	Long:  "Delete a resource by namespace, type and name",
	Run:   deleteHandler,
}

func deleteHandler(cmd *cobra.Command, args []string) {
	if len(args) != 3 {
		fmt.Println("Usage: delete <namespace> <resource type> <resource name>")
		os.Exit(1)
	}
	nameSpace := args[0]
	resourceType := args[1]
	resourceName := args[2]
	var url string

	switch resourceType {
	case "Pod":
		url = config.APIServerURL() + config.PodURI
	case "Service":
		url = config.APIServerURL() + config.ServiceURI
	case "ReplicaSet":
		url = config.APIServerURL() + config.ReplicaSetURI
	case "Dns":
		url = config.APIServerURL() + config.DNSURI
	default:
		fmt.Println("Supported resource types: Pod, Service, ReplicaSet, Dns")
	}

	url = strings.Replace(url, config.NameSpaceReplace, nameSpace, -1)
	url = strings.Replace(url, config.NameReplace, resourceName, -1)
	resp, err := httprequest.DelMsg(url, nil)
	if err != nil {
		fmt.Println("Error: Could not delete the object.")
		os.Exit(1)
	}

	DeleteResultDisplay(resourceName, resp)
}

func DeleteResultDisplay(name string, resp *http.Response) {
	if resp.StatusCode == 200 {
		fmt.Println(name + " deleted successfully.")
	} else {
		fmt.Println("Error: Could not delete the " + name + ".")
	}
}
