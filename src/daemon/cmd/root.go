// stock_ticker
// src/cmd/root.go

package cmd

import (
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"stock_ticker/i18n"
)

// langFlag holds the value of --lang; empty means auto-detect from the OS locale.
var langFlag string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "stock_ticker",
	Short:   "Add a short description here",
	Version: "0.0.1-1 (2026.08.10), Go version = " + runtime.Version(),
	Long: `This tools allows you to create a software directory structure.
This follows my template and allows you to package your software with minimal effort once built`,
	// Resolve the language once, before any subcommand runs.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return i18n.SetLanguage(langFlag)
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

	rootCmd.PersistentFlags().StringVarP(&langFlag, "lang", "l", "",
		"language ("+strings.Join(i18n.SupportedLanguages(), ", ")+"); default: OS locale")
}
