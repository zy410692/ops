package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/zy410692/ops/Lib"
)

func init() {
	rootCmd.AddCommand(SshCMD())
}

func SshCMD() *cobra.Command {
	var sshCMD = &cobra.Command{
		Use:   "ssh",
		Short: "ssh command ",
		Run: func(cmd *cobra.Command, args []string) {
			server := Lib.MustFlag("server", "string", cmd).(string)
			user := Lib.MustFlag("user", "string", cmd).(string)
			if session, err := Lib.SSHConnectKey(user, server, 22); err != nil {
				log.Fatal(err)
			} else {
				err := session.RequestPty("", 0, 0, Lib.ShellModes)
				if err != nil {
					log.Fatal(err)
				}
				session.Stdin = os.Stdin
				session.Stdout = os.Stdout
				session.Stderr = os.Stderr
				err = session.Run("/bin/bash")
				if err != nil {
					log.Fatal(err)
				}
			}

		},
	}
	sshCMD.Flags().StringP("server", "s", "", "set ssh server")
	sshCMD.Flags().StringP("user", "u", "", "set ssh username")
	return sshCMD

}
