package cmd

import (
	"fmt"
	httprequest "minik8s/internal/pkg/httpRequest"
	stringops "minik8s/internal/pkg/stringops"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
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
			url := config.APIServerUrl() + config.PodsURI
			nodes := []apiObject.Node{}
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
		case "pod":
			getPodHandler(namespace)
		case "service":
		case "deployment":
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
	url := config.APIServerUrl() + config.PodsURI
	url = strings.Replace(url, config.NameSpaceReplace, namespace, -1)
	pods := []apiObject.Pod{}
	code, err := httprequest.GetObjMsg(url, &pods, "data")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	if code.StatusCode != 200 {
		fmt.Println("Error: Failed to get pods")
		os.Exit(1)
	}
	printPodsResult(pods)
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
