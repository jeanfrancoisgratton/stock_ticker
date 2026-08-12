// stock_ticker
// src/cmd/root.go

package cmd

import (
	"fmt"
	"os"
	"runtime"

	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stock_ticker",
	Short: "Stock market simulator DAEMON",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the software version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(hftx.White("stockticker (daemon) v0.0.1 (2026.08.10), Go version = " + runtime.Version()))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.DisableAutoGenTag = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(versionCmd)
}
