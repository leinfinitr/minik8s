package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource by filename or stdin",
	Long:  "Create resources by filenames, stdin, resources and names, or by resources and label selector",
	Run:   createHandler,
}

func createHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: You must specify the type of resource to create.")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("Error: You may only specify one resource type to create.")
		os.Exit(1)
	}
	if len(args) == 1 {
		resourceType := args[0]
		if resourceType == "pod" {
			fmt.Println("Creating pod...")
		} else {
			fmt.Println("Error: You may only create pods. Try 'kubectl create --help' for more information.")
			os.Exit(1)
		}
	}
}