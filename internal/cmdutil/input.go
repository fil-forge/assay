// Package cmdutil holds helpers shared by assay commands.
package cmdutil

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// ReadInput reads the whole of the file named by the first argument, or of
// the command's stdin when no argument was given.
func ReadInput(cmd *cobra.Command, args []string) ([]byte, error) {
	if len(args) > 0 {
		b, err := os.ReadFile(args[0])
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", args[0], err)
		}
		return b, nil
	}
	b, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return nil, fmt.Errorf("reading stdin: %w", err)
	}
	if len(b) == 0 {
		return nil, fmt.Errorf("no input provided via stdin and no file specified")
	}
	return b, nil
}
