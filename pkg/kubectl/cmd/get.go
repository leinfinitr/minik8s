package cmd

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/httpRequest"
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

func init(){
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
			code, err := httprequest.GetObjMsg(url, &nodes, "data")
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
			if code.StatusCode != 200 {
				fmt.Println("Error: Failed to get nodes")
				os.Exit(1)
			}
			printNodesResult(nodes)
		}
		namespace, _ := cmd.Flags().GetString("namespace")
		if namespace == "" {
			namespace = "default"
		}
		switch resourceType {
		case apiObject.PodType:
			getPodHandler(namespace)
		case apiObject.ServiceType:
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
	resp,err := http.Get(url)
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

func getReplicaSetHandler(namespace string){
	url := config.APIServerURL() + config.ReplicaSetsURI
	url = strings.Replace(url, config.NameSpaceReplace, namespace, -1)
	var replicaSets []apiObject.ReplicaSet
	resp,err := http.Get(url)
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