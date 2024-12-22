package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ops",
	Short: "ops-运维命令行工具",
}

func Execute() {
	rootCmd.Execute()
}
