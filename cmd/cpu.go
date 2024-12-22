package cmd

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(CpuCMD())
}

func CpuCMD() *cobra.Command {
	var subCmd = &cobra.Command{
		Use:   "cpu",
		Short: "打印cpu百分比",
		Run: func(cmd *cobra.Command, args []string) {
			for {
				p, _ := cpu.Percent(3*time.Second, false)
				v, _ := mem.VirtualMemory()
				fmt.Printf("\rUsedPercent:%.1f%%---CPU:%.1f%%", v.UsedPercent, p[0])
				time.Sleep(time.Second)
			}
		},
	}
	return subCmd
}
