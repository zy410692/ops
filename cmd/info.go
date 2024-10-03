package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
	"os"
	"time"
)

func InfoCMD() *cobra.Command {
	var subCmd = &cobra.Command{
		Use:   "info",
		Short: "打印cpu内存百分比",
		Run: func(cmd *cobra.Command, args []string) {

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"项目", "数量", "百分比"})

			data := [][]string{}
			cpu_core, _ := cpu.Counts(true)
			cpu_percent, _ := cpu.Percent(time.Second, false)

			data = append(data, []string{"cpu", fmt.Sprintf("%d核", cpu_core), fmt.Sprintf("%1.f%%", cpu_percent[0])})
			m, _ := mem.VirtualMemory()
			data = append(data, []string{"memory", fmt.Sprintf("%dG", m.Total/1024/1024/1024),
				fmt.Sprintf("%1.f%%", m.UsedPercent)})

			for _, v := range data {
				table.Append(v)
			}

			table.Render()

		},
	}
	return subCmd
}
