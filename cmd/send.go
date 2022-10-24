package cmd

import (
	"github.com/abe1242/gosock/entity"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(sendCmd)
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send (or serve) a file",
	Long:  `Serve a file to those who request it`,
	Run: func(cmd *cobra.Command, args []string) {
		entity.Server(args[0], "0.0.0.0", "8888")
	},
}
