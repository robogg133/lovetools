package cmd

import (
	"path/filepath"

	"github.com/robogg133/lovetools/internal/build"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [directory]",
	Short: "build a love game",
	RunE: func(cmd *cobra.Command, args []string) error {
		builder, err := build.NewBuilder(true, true)
		if err != nil {
			return err
		}
		var dir string
		if len(args) > 0 && args[0] != "." {
			dir = args[0]
		} else {
			dir, err = filepath.Abs(".")
			if err != nil {
				return err
			}
		}
		return builder.Build(dir)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
