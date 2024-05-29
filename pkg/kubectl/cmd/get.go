package cmd

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/log"
	"minik8s/tools/stringops"
	"net/http"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a resource by filename or stdin",
	Long:  "Get resources by filenames, stdin, resources and names, or by resources and label selector",
	Run:   getHandler,
}

func init() {
	getCmd.PersistentFlags().StringP("namespace", "n", "", "Namespace")
}

func getHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: You must specify the type of resource to get.")
		os.Exit(1)
	}
	resourceType := args[0]
	println("resourceType:", resourceType)
	if stringops.StringInSlice(resourceType, apiObject.AllTypeList) {
		fmt.Println("Getting", resourceType)
	} else {
		fmt.Println("Error: The resource type specified is not supported.")
		os.Exit(1)
	}
	if len(args) == 1 {
		if resourceType == apiObject.NodeType {
			url := config.APIServerURL() + config.NodesURI
			var nodes []apiObject.Node
			res, err := http.Get(url)
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
			defer res.Body.Close()
			err = json.NewDecoder(res.Body).Decode(&nodes)
			if err != nil {
				log.ErrorLog("GetNodes: " + err.Error())
				os.Exit(1)
			}
			printNodesResult(nodes)
		}else if resourceType == apiObject.ContainerType{
			var pods []apiObject.Pod
			url := config.APIServerURL() + config.PodsURI
			res, err := http.Get(url)
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
			defer res.Body.Close()
			err = json.NewDecoder(res.Body).Decode(&pods)
			if err != nil {
				log.ErrorLog("GetPod: " + err.Error())
				os.Exit(1)
			}
			var containers []apiObject.Container
			for _, pod := range pods {
				containers = append(containers, pod.Spec.Containers...)
			}
			printContainersResult(containers)
		}
		namespace, _ := cmd.Flags().GetString("namespace")
		if namespace == "" {
			namespace = "default"
		}
		switch resourceType {
		case apiObject.PodType:
			getPodHandler(namespace)
		case apiObject.ServiceType:
			getServiceHandler(namespace)
		case apiObject.ReplicaSetType:
			getReplicaSetHandler(namespace)
		}
	}
}

func printNodesResult(nodes []apiObject.Node) {
	writer := table.NewWriter()
	writer.SetOutputMirror(os.Stdout)
	writer.AppendHeader(table.Row{"Kind", "Name", "Status", "IP", "Role"})
	for _, node := range nodes {
		printNodeResult(node, writer)
	}
	writer.Render()
}
func printContainersResult(containers []apiObject.Container) {
	writer := table.NewWriter()
	writer.SetOutputMirror(os.Stdout)
	writer.AppendHeader(table.Row{"Kind", "Name", "Image", "Status", "Ports"})
	for _, container := range containers {
		printContainerResult(container, writer)
	}
	writer.Render()
}
func printContainerResult(container apiObject.Container, writer table.Writer) {
	// 根据状态为Status单元格选择颜色
	var statusColor text.Colors
	switch container.ContainerStatus {
	case apiObject.ContainerRunning:
		statusColor = text.Colors{text.FgGreen}
	case apiObject.ContainerUnknown:
		statusColor = text.Colors{text.FgYellow}
	case apiObject.ContainerExited:
		statusColor = text.Colors{text.FgRed}
	default:
		statusColor = text.Colors{text.FgWhite}
	}

	// 应用颜色到Status
	coloredStatus := statusColor.Sprint(container.ContainerStatus)

	writer.AppendRow(table.Row{
		"Container",
		container.Name,
		container.Image,
		coloredStatus, // 使用包装了颜色的status
		container.Ports,
	})
}
func printNodeResult(node apiObject.Node, writer table.Writer) {
	// 根据状态为Status单元格选择颜色
	var statusColor text.Colors
	switch node.Status.Conditions[0].Type {
	case "Ready":
		statusColor = text.Colors{text.FgGreen}
	default:
		statusColor = text.Colors{text.FgRed}
	}

	// 应用颜色到Status
	coloredStatus := statusColor.Sprint(node.Status.Conditions[0].Type)

	var roleColor text.Colors
	switch node.Metadata.Labels["kubernetes.io/role"] {
	case "master":
		roleColor = text.Colors{text.FgHiYellow, text.BgBlack}
	case "worker":
		roleColor = text.Colors{text.FgHiBlue, text.BgBlack}
	default:
		roleColor = text.Colors{text.FgWhite, text.BgBlack}
	}

	// 应用颜色到Role
	coloredRole := roleColor.Sprint(node.Metadata.Labels["kubernetes.io/role"])

	writer.AppendRow(table.Row{
		"Node",
		node.Metadata.Name,
		coloredStatus, // 使用包装了颜色的status
		node.Status.Addresses[0].Address,
		coloredRole, // 使用包装了颜色的role
	})
}

func getPodHandler(namespace string) {
	url := config.APIServerURL() + config.PodsURI
	url = strings.Replace(url, config.NameSpaceReplace, namespace, -1)
	var pods []apiObject.Pod
	resp, err := http.Get(url)
	if err != nil {
		log.ErrorLog("GetPod: " + err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&pods)
	if err != nil {
		log.ErrorLog("GetPod: " + err.Error())
		os.Exit(1)
	}
	printPodsResult(pods)
}

func getServiceHandler(namespace string){
	url := config.APIServerURL() + config.ServicesURI
	url = strings.Replace(url, config.NameSpaceReplace, namespace, -1)
	var services []apiObject.Service
	resp, err := http.Get(url)
	if err != nil {
		log.ErrorLog("GetService: " + err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&services)
	if err != nil {
		log.ErrorLog("GetService: " + err.Error())
		os.Exit(1)
	}
	printServicesResult(services)

}
func getReplicaSetHandler(namespace string) {
	url := config.APIServerURL() + config.ReplicaSetsURI
	url = strings.Replace(url, config.NameSpaceReplace, namespace, -1)
	var replicaSets []apiObject.ReplicaSet
	resp, err := http.Get(url)
	if err != nil {
		log.ErrorLog("GetReplicaSet: " + err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&replicaSets)
	if err != nil {
		log.ErrorLog("GetReplicaSet: " + err.Error())
		os.Exit(1)
	}
	printReplicasetsResult(replicaSets)
}

func printServicesResult(services []apiObject.Service) {
	writer := table.NewWriter()
	writer.SetOutputMirror(os.Stdout)
	writer.AppendHeader(table.Row{"Kind", "Name", "ClusterIP", "Ports"})
	for _, service := range services {
		printServiceResult(service, writer)
	}
	writer.Render()
}

func printServiceResult(service apiObject.Service, writer table.Writer) {
	ports := ""
	for _, port := range service.Spec.Ports {
		ports += fmt.Sprintf("%d/%s ", port.Port, port.Protocol)
	}
	writer.AppendRow(table.Row{
		"Service",
		service.Metadata.Name,
		service.Spec.ClusterIP,
		ports,
	})
}

func printPodsResult(pods []apiObject.Pod) {
	writer := table.NewWriter()
	writer.SetOutputMirror(os.Stdout)
	writer.AppendHeader(table.Row{"Kind", "Name", "Status", "Node"})
	for _, pod := range pods {
		printPodResult(pod, writer)
	}
	writer.Render()
}

func printPodResult(pod apiObject.Pod, writer table.Writer) {
	// 根据状态为Status单元格选择颜色
	var statusColor text.Colors
	switch pod.Status.Phase {
	case "Running":
		statusColor = text.Colors{text.FgGreen}
	case "Pending":
		statusColor = text.Colors{text.FgYellow}
	case "Failed":
		statusColor = text.Colors{text.FgRed}
	default:
		statusColor = text.Colors{text.FgWhite}
	}

	// 应用颜色到Status
	coloredStatus := statusColor.Sprint(pod.Status.Phase)

	writer.AppendRow(table.Row{
		"Pod",
		pod.Metadata.Name,
		coloredStatus, // 使用包装了颜色的status
		pod.Spec.NodeName,
	})
}

func printReplicasetsResult(replicaSets []apiObject.ReplicaSet) {
	writer := table.NewWriter()
	writer.SetOutputMirror(os.Stdout)
	writer.AppendHeader(table.Row{"Kind", "Name", "Replicas", "Selector"})
	for _, rs := range replicaSets {
		printReplicaSetResult(rs, writer)
	}
	writer.Render()
}

func printReplicaSetResult(rs apiObject.ReplicaSet, writer table.Writer) {
	writer.AppendRow(table.Row{
		"ReplicaSet",
		rs.Metadata.Name,
		rs.Spec.Replicas,
		rs.Spec.Selector,
	})
}
