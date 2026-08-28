package cid

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:          "parse <cid>",
	Short:        "Print information about a string encoded CID.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("parsing CID: %w", err)
		}
		return printCID(cmd.OutOrStdout(), c)
	},
}

func init() {
	Cmd.AddCommand(parseCmd)
}
