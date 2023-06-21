package cmd

import (
	"github.com/abe1242/gosock/entity"

	"github.com/spf13/cobra"
)

var sclip bool

func init() {
	sendCmd.Flags().BoolVarP(&sclip, "clipboard", "k", false, "set this flag to send clipboard")

	rootCmd.AddCommand(sendCmd)
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send (or serve) a file",
	Long:  `Serve a file to those who request it`,
	Run: func(cmd *cobra.Command, args []string) {
		var file string
		if len(args) > 0 {
			file = args[0]
		} else {
			file = ""
		}
		entity.Server(file, "0.0.0.0", "8888", sclip)
	},
}
