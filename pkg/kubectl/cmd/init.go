package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "Kubernetes CLI",
	Long:  "Kubernetes CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Kubernetes CLI")
		fmt.Println(cmd.UsageString())
	},
}

func init() {
	rootCmd.AddCommand(deletedCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(describeCmd)
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(getCmd)
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
