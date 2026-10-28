package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Love version tools",
	Long:  "Love version tools and show",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("1.0.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
