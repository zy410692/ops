package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCmd 创建根命令
func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "myapp",
		Short: "MyApp is a CLI application",
	}
	return rootCmd
}
