package cmd

import (
	"fmt"

	"github.com/zy410692/ops/Lib"

	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "生成密码相关命令",
	Long:  `生成密码的相关命令，包括创建新密码和列出已有密码`,
}

var createCmd = &cobra.Command{
	Use:   "create [网站/信息] [密码]",
	Short: "为指定网站生成或设置新密码",
	Long:  `为指定的网站或信息生成一个新的8位随机密码，或使用指定的密码`,
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		site := args[0]
		pm := Lib.NewPasswordManager()

		if err := pm.LoadFromFile(); err != nil {
			return fmt.Errorf("加载密码文件失败: %v", err)
		}

		var err error
		if len(args) == 2 {
			// 使用指定的密码
			err = pm.AddPassword(site, args[1])
		} else {
			// 自动生成密码
			err = pm.AddPassword(site)
		}

		if err != nil {
			return fmt.Errorf("生成密码失败: %v", err)
		}

		fmt.Printf("已为 %s 生成新密码\n", site)

		// 自动显示密码列表
		entries := pm.ListPasswords()
		if len(entries) == 0 {
			fmt.Println("没有保存的密码")
			return nil
		}

		fmt.Println("\n当前保存的密码列表:")
		for _, entry := range entries {
			fmt.Printf("网站/信息: %s, 密码: %s\n", entry.Site, entry.Password)
		}
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有保存的密码",
	Long:  `显示所有已保存的网站和对应的密码`,
	RunE: func(cmd *cobra.Command, args []string) error {
		pm := Lib.NewPasswordManager()

		if err := pm.LoadFromFile(); err != nil {
			return fmt.Errorf("加载密码文件失败: %v", err)
		}

		entries := pm.ListPasswords()
		if len(entries) == 0 {
			fmt.Println("没有保存的密码")
			return nil
		}

		fmt.Println("保存的密码列表:")
		for _, entry := range entries {
			fmt.Printf("网站/信息: %s, 密码: %s\n", entry.Site, entry.Password)
		}
		return nil
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [网站/信息|all]",
	Short: "删除指定网站的密码或删除所有密码",
	Long:  `删除指定网站或信息对应的密码记录，使用 'all' 参数可以删除所有密码`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		site := args[0]
		pm := Lib.NewPasswordManager()

		if err := pm.LoadFromFile(); err != nil {
			return fmt.Errorf("加载密码文件失败: %v", err)
		}

		if site == "all" {
			if err := pm.DeleteAllPasswords(); err != nil {
				return fmt.Errorf("删除所有密码失败: %v", err)
			}
			fmt.Println("已删除所有密码")
		} else {
			if err := pm.DeletePassword(site); err != nil {
				return fmt.Errorf("删除密码失败: %v", err)
			}
			fmt.Printf("已删除 %s 的密码\n", site)
		}

		// 自动显示更新后的密码列表
		entries := pm.ListPasswords()
		if len(entries) == 0 {
			fmt.Println("\n当前没有保存的密码")
			return nil
		}

		fmt.Println("\n当前保存的密码列表:")
		for _, entry := range entries {
			fmt.Printf("网站/信息: %s, 密码: %s\n", entry.Site, entry.Password)
		}
		return nil
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [网站/信息] [密码]",
	Short: "更新指定网站的密码",
	Long:  `为指定的网站或信息更新密码，可以自动生成新密码或使用指定的密码`,
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		site := args[0]
		pm := Lib.NewPasswordManager()

		if err := pm.LoadFromFile(); err != nil {
			return fmt.Errorf("加载密码文件失败: %v", err)
		}

		var err error
		if len(args) == 2 {
			// 使用指定的密码
			err = pm.UpdatePassword(site, args[1])
		} else {
			// 自动生成密码
			err = pm.UpdatePassword(site)
		}

		if err != nil {
			return fmt.Errorf("更新密码失败: %v", err)
		}

		fmt.Printf("已更新 %s 的密码\n", site)
		return nil
	},
}

var verifyCmd = &cobra.Command{
	Use:   "verify [网站/信息]",
	Short: "验证指定网站是否存在密码",
	Long:  `检查指定的网站或信息是否已经存在密码记录`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		site := args[0]
		pm := Lib.NewPasswordManager()

		if err := pm.LoadFromFile(); err != nil {
			return fmt.Errorf("加载密码文件失败: %v", err)
		}

		exists, entry := pm.VerifyPassword(site)
		if exists {
			fmt.Printf("网站 %s 的密码信息:\n", site)
			fmt.Printf("密码: %s\n", entry.Password)
			fmt.Printf("创建时间: %s\n", entry.Created.Format("2006-01-02 15:04:05"))
			fmt.Printf("最后修改: %s\n", entry.Modified.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("网站 %s 没有保存密码\n", site)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.AddCommand(createCmd)
	genCmd.AddCommand(listCmd)
	genCmd.AddCommand(deleteCmd)
	genCmd.AddCommand(updateCmd)
	genCmd.AddCommand(verifyCmd)
}
