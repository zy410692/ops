package cmd

import (
	"context"
	"fmt"
	"log"

	"context"
	"fmt"
	"log"

	"github.com/mittwald/goharbor-client/v4/apiv2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			fmt.Println("子命令行逻辑")
			//初始化viper
			viper.SetConfigName("config") // 配置文件名称（不包括扩展名）
			viper.SetConfigType("yaml")   // 配置文件类型
			viper.AddConfigPath(".")      // 配置文件搜索路径

			// 读取配置文件
			err := viper.ReadInConfig()
			if err != nil {
				fmt.Println("错误：读取配置文件失败", err)
				return
			}

			// 获取环境变量
			harbor_url := viper.GetString("harbor.harbor_url")
			harbor_user := viper.GetString("harbor.user")
			harbor_user_password := viper.GetString("harbor.password")
			harbor_project_name := viper.GetString("harbor.harbor_project_name")

			if harbor_url == "" {
				panic("harbor url 参数不能为空")
			}
			if harbor_user == "" {
				panic("config.yaml harbor root user 参数不能为空")
			}
			if harbor_user_password == "" {
				panic(" harbor user 参数不能为空")
			}
			if harbor_project_name == "" {
				panic("config.yaml harbor project name 参数不能为空")
			}

			harborClient, err := apiv2.NewRESTClientForHost(fmt.Sprintf("http://%s/api/", harbor_url), harbor_user, harbor_user_password)
			if err != nil {
				log.Println("连接harbor错误", err)
			}

			project, _ := harborClient.GetProject(context.Background(), harbor_project_name)

			if project == nil {
				var storageLimit int64 = -1
				result2, err := harborClient.NewProject(context.TODO(), harbor_project_name, &storageLimit)
				if err != nil {
					log.Println("创建项目错误", err)
				}
				fmt.Println(result2)
			}

		},
	}

	return subCmd
}
