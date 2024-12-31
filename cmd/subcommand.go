package cmd

import (
	"fmt"
	"log"

	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zy410692/ops/Lib"
)

func init() {

	rootCmd.AddCommand(NewSubcommand())
}

// NewSubcommand 创建一个新的子命令
func NewSubcommand() *cobra.Command {
	var subCmd = &cobra.Command{
		Use:   "minio",
		Short: "minio添加用户和鉴权",
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
			minio_url := viper.GetString("minio.minio_ip")
			minio_port := viper.GetInt("minio.minio_port")
			minio_root_user := viper.GetString("minio.minio_user")
			minio_password := viper.GetString("minio.minio_password")

			//获取参数
			minio_bucket, _ := cmd.Flags().GetString("minio.bucket")
			minio_user, _ := cmd.Flags().GetString("minio. suser")

			if minio_url == "" {
				panic("minio url 参数不能为空")
			}
			if minio_root_user == "" {
				panic("config.yaml minio root user 参数不能为空")
			}
			if minio_user == "" {
				panic(" minio user 参数不能为空")
			}
			if minio_password == "" {
				panic("config.yaml minio password 参数不能为空")
			}
			if minio_bucket == "" {
				panic("minio bucket 参数不能为空")
			}

			Lib.Initialize(minio_url+":"+strconv.Itoa(minio_port), minio_root_user, minio_password)
			Lib.InitMinio(minio_url+":"+strconv.Itoa(minio_port), minio_root_user, minio_password)

			password := Lib.GeneratePassword(8)
			log.Println("password:***********")
			log.Println(password)
			log.Println("--------:***********")
			mdmclient := Lib.GetMinioAdmClient()
			if Lib.ExistUser(mdmclient, minio_user) {
				log.Printf("%s minio账户已经存在", minio_user)
			} else {
				log.Printf("%s minio账户不存在 创建账户", minio_user)
				Lib.CreateUser(mdmclient, minio_user, password)
			}

			log.Printf("检测到newbucket名称 %s", minio_bucket)
			Lib.CreateBucket(minio_bucket)
			//创建policy--policy绑定bucketname&user
			Lib.SetPolicySingle(mdmclient, minio_bucket, minio_user)
			//helpers.SetPolicyMult(mdmclient, minio_bucket, minio_user)

		},
	}

	// 添加子命令的参数（如果有）
	subCmd.Flags().StringP("user", "u", "", "minio的user")
	subCmd.Flags().StringP("bucket", "b", "", "minio的bucketname")
	return subCmd
}
