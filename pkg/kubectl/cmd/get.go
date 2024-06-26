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
		namespace, _ := cmd.Flags().GetString("namespace")
		if namespace == "" {
			namespace = "default"
		}
		switch resourceType {
		case apiObject.NodeType:
			getNodeHandler()
		case apiObject.ContainerType:
			getContainerHandler()
		case apiObject.PodType:
			getPodHandler(namespace)
		case apiObject.ServiceType:
			getServiceHandler(namespace)
		case apiObject.ReplicaSetType:
			getReplicaSetHandler(namespace)
		case apiObject.HpaType:
			getHpaHandler(namespace)
		}
	}
}

func getNodeHandler() {
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

func getContainerHandler() {
	var pods []apiObject.Pod
	url := config.APIServerURL() + config.PodsGlobalURI
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

func printContainersResult(containers []apiObject.Container) {
	writer := table.NewWriter()
	writer.SetOutputMirror(os.Stdout)
	writer.AppendHeader(table.Row{"Kind", "Name", "ID", "Image", "Status", "Ports"})
	for _, container := range containers {
		printContainerResult(container, writer)
	}
	writer.Render()
}

func printContainerResult(container apiObject.Container, writer table.Writer) {
	// 根据状态为Status单元格选择颜色
	var statusColor text.Colors
	switch container.ContainerStatus {
	case apiObject.ContainerCreated:
		statusColor = text.Colors{text.FgGreen}
	case apiObject.ContainerRunning:
		statusColor = text.Colors{text.FgGreen}
	case apiObject.ContainerExited:
		statusColor = text.Colors{text.FgRed}
	case apiObject.ContainerUncreated:
		statusColor = text.Colors{text.FgYellow}
	case apiObject.ContainerUnknown:
		statusColor = text.Colors{text.FgWhite}
	default:
		statusColor = text.Colors{text.FgWhite}
	}
	var containerStatusMap = map[apiObject.ContainerStatus]string{
		apiObject.ContainerCreated:   "Created",
		apiObject.ContainerRunning:   "Running",
		apiObject.ContainerExited:    "Exited",
		apiObject.ContainerUncreated: "Uncreated",
		apiObject.ContainerUnknown:   "Unknown",
	}
	// 应用颜色到Status
	coloredStatus := statusColor.Sprint(containerStatusMap[container.ContainerStatus])

	writer.AppendRow(table.Row{
		"Container",
		container.Name,
		container.ContainerID,
		container.Image,
		coloredStatus, // 使用包装了颜色的status
		container.Ports,
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

func printPodsResult(pods []apiObject.Pod) {
	writer := table.NewWriter()
	writer.SetOutputMirror(os.Stdout)
	writer.AppendHeader(table.Row{"Kind", "Namespace", "Name", "Status", "Node", "IP"})
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
		pod.Metadata.Namespace,
		pod.Metadata.Name,
		coloredStatus, // 使用包装了颜色的status
		pod.Spec.NodeName,
		pod.Status.PodIP,
	})
}

func getServiceHandler(namespace string) {
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

func printServicesResult(services []apiObject.Service) {
	writer := table.NewWriter()
	writer.SetOutputMirror(os.Stdout)
	writer.AppendHeader(table.Row{"Kind", "Namespace", "Name", "ClusterIP", "Ports"})
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
		service.Metadata.Namespace,
		service.Metadata.Name,
		service.Spec.ClusterIP,
		ports,
	})
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

func printReplicasetsResult(replicaSets []apiObject.ReplicaSet) {
	writer := table.NewWriter()
	writer.SetOutputMirror(os.Stdout)
	writer.AppendHeader(table.Row{"Kind", "Namespace", "Name", "Replicas", "Selector"})
	for _, rs := range replicaSets {
		printReplicaSetResult(rs, writer)
	}
	writer.Render()
}

func printReplicaSetResult(rs apiObject.ReplicaSet, writer table.Writer) {
	writer.AppendRow(table.Row{
		"ReplicaSet",
		rs.Metadata.Namespace,
		rs.Metadata.Name,
		rs.Spec.Replicas,
		rs.Spec.Selector,
	})
}

func getHpaHandler(namespace string) {
	url := config.APIServerURL() + config.HpasURI
	url = strings.Replace(url, config.NameSpaceReplace, namespace, -1)
	var hpas []apiObject.HPA
	resp, err := http.Get(url)
	if err != nil {
		log.ErrorLog("GetHpa: " + err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.ErrorLog("GetHpa: " + resp.Status)
		os.Exit(1)
	}
	err = json.NewDecoder(resp.Body).Decode(&hpas)
	if err != nil {
		log.ErrorLog("GetHpa: " + err.Error())
		os.Exit(1)
	}
	printHpasResult(hpas)
}

func printHpasResult(hpas []apiObject.HPA) {
	writer := table.NewWriter()
	writer.SetOutputMirror(os.Stdout)
	writer.AppendHeader(table.Row{"Kind", "Namespace", "Name", "MinReplicas", "MaxReplicas", "CurrentReplicas"})
	for _, hpa := range hpas {
		printHpaResult(hpa, writer)
	}
	writer.Render()
}

func printHpaResult(hpa apiObject.HPA, writer table.Writer) {
	writer.AppendRow(table.Row{
		"Hpa",
		hpa.Metadata.Namespace,
		hpa.Metadata.Name,
		hpa.Spec.MinReplicas,
		hpa.Spec.MaxReplicas,
		hpa.Status.CurrentReplicas,
	})
}
