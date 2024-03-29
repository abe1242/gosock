package cmd

import (
	"github.com/abe1242/gosock/entity"

	"github.com/spf13/cobra"
)

var (
	contnue bool
	outfile string
	rclip   bool
)

func init() {
	recvCmd.Flags().BoolVarP(&contnue, "continue", "c", false, "set this flag to resume downloads")
	recvCmd.Flags().BoolVarP(&rclip, "clipboard", "k", false, "set this flag to copy output to clipboard")
	recvCmd.Flags().StringVarP(&outfile, "output", "o", "", "specify the output filename")

	rootCmd.AddCommand(recvCmd)
}

var recvCmd = &cobra.Command{
	Use:   "recv",
	Short: "Recieve (or request) a file",
	Long:  `Recieve a file from a server`,
	Run: func(cmd *cobra.Command, args []string) {
		entity.Client(args[0], "8888", contnue, rclip, outfile)
	},
}
