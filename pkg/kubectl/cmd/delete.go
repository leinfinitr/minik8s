package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var deletedCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a resource by filename or stdin",
	Long:  "Delete resources by filenames, stdin, resources and names, or by resources and label selector",
	Run:   deleteHandler,
}

func deleteHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: You must specify the type of resource to delete.")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("Error: You may only specify one resource type to delete.")
		os.Exit(1)
	}
	if len(args) == 1 {
		resourceType := args[0]
		if resourceType == "pod" {
			fmt.Println("Deleting pod...")
		} else {
			fmt.Println("Error: You may only delete pods. Try 'kubectl delete --help' for more information.")
			os.Exit(1)
		}
	}
}
