package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a configuration to a resource by filename or stdin",
	Long:  "Apply a configuration to a resource by filename or stdin",
	Run:   applyHandler,
}

func applyHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: You must specify the type of resource to apply.")
		os.Exit(1)
	}
	if len(args) > 1 {
		fmt.Println("Error: You may only specify one resource type to apply.")
		os.Exit(1)
	}
	if len(args) == 1 {
		resourceType := args[0]
		if resourceType == "pod" {
			fmt.Println("Applying pod...")
		} else {
			fmt.Println("Error: You may only apply pods. Try 'kubectl apply --help' for more information.")
			os.Exit(1)
		}
	}
}