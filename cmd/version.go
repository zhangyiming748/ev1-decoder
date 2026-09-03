package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd is the "ev1 version" subcommand. It prints the same version
// string that backs the root command's -v/--version flag.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the ev1 version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ev1 version %s\n", rootCmd.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
