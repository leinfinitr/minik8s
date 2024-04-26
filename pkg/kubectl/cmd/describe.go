package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe a resource by filename or stdin",
	Long:  "Describe resources by filenames, stdin, resources and names, or by resources and label selector",
	Run:   describeHandler,
}

func describeHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: You must specify the type of resource to describe.")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("Error: You may only specify one resource type to describe.")
		os.Exit(1)
	}
	if len(args) == 1 {
		resourceType := args[0]
		if resourceType == "pod" {
			fmt.Println("Describing pod...")
		} else {
			fmt.Println("Error: You may only describe pods. Try 'kubectl describe --help' for more information.")
			os.Exit(1)
		}
	}
}
