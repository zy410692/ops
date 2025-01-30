package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/mittwald/goharbor-client/v4/apiv2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configFileName = "config" // 配置文件名称
	configFileType = "yaml"   // 配置文件类型
)

func init() {
	rootCmd.AddCommand(HarborCMD())
}

// NewSubcommand 创建一个新的子命令
func HarborCMD() *cobra.Command {
	var subCmd = &cobra.Command{
		Use:   "harbor",
		Short: "harbor 命令行工具",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runHarborCommand(); err != nil {
				log.Fatalf("执行命令时出错: %v", err)
			}
		},
	}

	return subCmd
}

// runHarborCommand 执行 harbor 命令的逻辑
func runHarborCommand() error {
	// 初始化 viper
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(".")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 获取环境变量
	harborURL := viper.GetString("harbor.harbor_url")
	harborUser := viper.GetString("harbor.user")
	harborUserPassword := viper.GetString("harbor.password")
	harborProjectName := viper.GetString("harbor.harbor_project_name")

	// 参数检查
	if harborURL == "" {
		return fmt.Errorf("harbor url 参数不能为空")
	}
	if harborUser == "" {
		return fmt.Errorf("config.yaml harbor root user 参数不能为空")
	}
	if harborUserPassword == "" {
		return fmt.Errorf("harbor user 参数不能为空")
	}
	if harborProjectName == "" {
		return fmt.Errorf("config.yaml harbor project name 参数不能为空")
	}

	// 创建 Harbor 客户端
	harborClient, err := apiv2.NewRESTClientForHost(fmt.Sprintf("http://%s/api/", harborURL), harborUser, harborUserPassword)
	if err != nil {
		return fmt.Errorf("连接 harbor 错误: %w", err)
	}
	defer func() {

	}()

	ctx := context.Background()
	project, err := harborClient.GetProject(ctx, harborProjectName)
	if err != nil {
		return fmt.Errorf("获取项目错误: %w", err)
	}

	if project == nil {
		var storageLimit int64 = -1
		result2, err := harborClient.NewProject(ctx, harborProjectName, &storageLimit)
		if err != nil {
			return fmt.Errorf("创建项目错误: %w", err)
		}
		fmt.Println(result2)
	}

	return nil
}
