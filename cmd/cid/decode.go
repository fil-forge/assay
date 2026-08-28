package cid

import (
	"fmt"
	"os"

	"github.com/ipfs/go-cid"
	"github.com/spf13/cobra"
)

var decodeCmd = &cobra.Command{
	Use:          "decode <path>",
	Short:        "Print information about a byte encoded CID.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("reading %s: %w", args[0], err)
		}
		c, err := cid.Cast(b)
		if err != nil {
			return fmt.Errorf("decoding CID: %w", err)
		}
		return printCID(cmd.OutOrStdout(), c)
	},
}

func init() {
	Cmd.AddCommand(decodeCmd)
}
