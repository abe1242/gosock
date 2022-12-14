package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gosock",
	Short: "send files over the network",
	Long: `Send your files over the network blazingly fast!
It uses a custom protocol on top of TCP for filesharing`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
