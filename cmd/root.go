package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "ops",
	Short: "ops-运维命令行工具",
	Run: func(cmd *cobra.Command, args []string) {
		// 读取配置文件
		if err := loadConfig(configFile); err != nil {
			log.Fatalf("加载配置文件失败: %v", err)
		}
	},
}

func init() {
	// 添加标志到父命令
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "config.yaml", "配置文件路径")
}

func loadConfig(filePath string) error {
	// 这里可以实现加载 YAML 配置文件的逻辑
	// 例如使用 viper 库来读取 YAML 文件
	// viper.SetConfigFile(filePath)
	// return viper.ReadInConfig()

	// 示例：简单打印文件路径
	log.Printf("加载配置文件: %s", filePath)
	return nil
}

func Execute() {
	rootCmd.Execute()
}
