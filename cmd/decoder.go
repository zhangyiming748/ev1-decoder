package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/zhangyiming748/ev1-decoder/core"
)

// dir is the target directory scanned by the decoder subcommand.
var dir string

// decoderCmd is the "ev1 decoder" subcommand. It walks a directory and
// decodes every *.ev1 file it finds using the core package.
var decoderCmd = &cobra.Command{
	Use:   "decoder",
	Short: "Decode all .ev1 files in a directory into .flv",
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := os.Stat(dir)
		if err != nil {
			return fmt.Errorf("cannot access %q: %w", dir, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("%q is not a directory", dir)
		}

		var decoded, failed int
		walkErr := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				// Skip entries we cannot read but keep scanning.
				fmt.Fprintf(cmd.ErrOrStderr(), "[skip] %s: %v\n", path, err)
				return nil
			}
			if d.IsDir() || filepath.Ext(path) != core.EV1Ext {
				return nil
			}

			newPath, derr := core.DecryptFile(path)
			if derr != nil {
				failed++
				fmt.Fprintf(cmd.ErrOrStderr(), "[fail] %s: %v\n", path, derr)
				return nil
			}
			decoded++
			fmt.Fprintf(cmd.OutOrStdout(), "[ok]   %s -> %s\n", path, newPath)
			return nil
		})
		if walkErr != nil {
			return walkErr
		}

		fmt.Fprintf(cmd.OutOrStdout(), "\nDone. %d decoded, %d failed.\n", decoded, failed)
		return nil
	},
}

func init() {
	decoderCmd.Flags().StringVarP(&dir, "dir", "d", "", "directory to scan for .ev1 files")
	_ = decoderCmd.MarkFlagRequired("dir")
	rootCmd.AddCommand(decoderCmd)
}
