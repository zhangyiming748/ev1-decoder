// Package cmd wires up the ev1 command line interface using cobra.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the top level "ev1" command.
var rootCmd = &cobra.Command{
	Use:   "ev1",
	Short: "EV1 video format toolkit",
	Long: "ev1 is a small toolkit for the EV1 video format.\n" +
		"An EV1 file is a regular FLV whose first 100 bytes are masked with XOR 0xFF.",
	// A bare "ev1" invocation just prints help.
	SilenceUsage: true,
}

// Execute runs the root command and exits non-zero on error.
// The version string comes from main (injected via -ldflags, see .goreleaser.yml)
// and backs both the root -v/--version flag and the "version" subcommand.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
