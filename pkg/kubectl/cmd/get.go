package cmd

import (
	"fmt"
	"os"

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
	if len(args) > 1 {
		fmt.Println("Error: You may only specify one resource type to get.")
		os.Exit(1)
	}
	if len(args) == 1 {
		resourceType := args[0]
		if resourceType == "pod" {
			fmt.Println("Getting pod...")
		} else {
			fmt.Println("Error: You may only get pods. Try 'kubectl get --help' for more information.")
			os.Exit(1)
		}
	}
}
