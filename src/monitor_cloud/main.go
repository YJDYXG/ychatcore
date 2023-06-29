package main

import (
	"ychatcore"
	"fmt"
	"github.com/spf13/cobra"
)

func main() {
	var version bool
	var proto string
	var addr string
	var name string
	var conf string

	var rootCmd = &cobra.Command{
		Use:	"systemd",
		Short:	"managed by Systemd",
		Run: func(cmd *cobra.Command, args []string) {
			if version {
				fmt.Println(VERSION)
			}
		},
	}

	var ychatcoreCmd = &cobra.Command{
		Use:   "ychatcore",
		Short: "ychatcore is the backend core process of chat ai",
		Run: func(cmd *cobra.Command, args []string) {
			if version {
				fmt.Println(VERSION)
			}
		},
	}

	var ychatcoreStartCmd = &cobra.Command{ //正常命令
		Use:   "start",
		Short: "start ychatcore daemon",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("start ychatcore daemon")
			a, err := ychatcore.New(conf, name, proto, addr)
			if err != nil {
				log.Error("new service error: %v", err)
			}
			a.Start()
		},
	}

	ychatcoreStartCmd.Flags().StringVarP(&conf, "conf", "c", "/var/lib/ychatcore/ychatcore.conf", "service config file")
	ychatcoreStartCmd.Flags().StringVarP(&name, "name", "n", "ychatcore", "service's name")
	ychatcoreStartCmd.Flags().StringVarP(&proto, "proto", "p", "", "service protocol")
	ychatcoreStartCmd.Flags().StringVarP(&addr, "addr", "a", "", "service listen addr")

	ychatcoreCmd.AddCommand(ychatcoreStartCmd)
	rootCmd.AddCommand(ychatcoreCmd)
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "ychatcore version") //参数flag
	rootCmd.Execute()
}