package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"main/cmd"
)

func main() {
	// 创建根命令
	var rootCmd = &cobra.Command{
		Use:   "ops",
		Short: "ops-运维命令行工具",
	}

	// 添加子命令
	rootCmd.AddCommand(cmd.NewSubcommand())
	rootCmd.AddCommand(cmd.CpuCMD())
	rootCmd.AddCommand(cmd.InfoCMD())
	rootCmd.AddCommand(cmd.SshCMD())

	// 执行根命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
