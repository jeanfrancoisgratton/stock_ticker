// stock_ticker
// src/cmd/root.go

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "stock_ticker",
	Short:   "Add a short description here",
	Version: "0.0.1-1 (2026.08.10), Go version = " + runtime.Version(),
	Long: `This tools allows you to create a software directory structure.
This follows my template and allows you to package your software with minimal effort once built`,
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
}
