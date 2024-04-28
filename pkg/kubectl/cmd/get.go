package cmd

import (
	"fmt"
	httprequest "minik8s/internal/pkg/httpRequest"
	stringops "minik8s/internal/pkg/stringops"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"os"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jedib0t/go-pretty/v6/table"
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
	if stringops.StringInSlice(resourceType, apiObject.AllTypeList) {
		fmt.Println("Getting", resourceType)
	} else {
		fmt.Println("Error: The resource type specified is not supported.")
		os.Exit(1)
	}
	if len(args) == 1 {
		if resourceType == "Node" {
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

	// 应用颜色到Role，这里只是一个例子，实际应用时你需要根据实际情况来判断
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