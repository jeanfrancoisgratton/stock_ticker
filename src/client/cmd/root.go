// stock_ticker
// src/cmd/root.go

package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"stock_ticker/i18n"

	hftx "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
)

// langFlag holds the value of --lang; empty means auto-detect from the OS locale.
var langFlag string

var rootCmd = &cobra.Command{
	Use:   "stock_ticker",
	Short: "Stock market simulator CLIENT",
	// Resolve the language once, before any subcommand runs.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return i18n.SetLanguage(langFlag)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the software version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(hftx.White("stockticker v0.0.1-1 (2026.08.10), Go version = " + runtime.Version()))
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

	rootCmd.PersistentFlags().StringVarP(&langFlag, "lang", "l", "",
		"language ("+strings.Join(i18n.SupportedLanguages(), ", ")+"); default: OS locale")
}
