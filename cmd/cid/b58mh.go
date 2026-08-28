package cid

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multibase"
	"github.com/spf13/cobra"
)

var b58mhCmd = &cobra.Command{
	Use:          "b58mh <cid>",
	Short:        "Print the multibase base58btc encoded multihash from a CID.",
	Long:         "Extract the multihash from a CID and print the multibase base58btc encoded string.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := cid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("parsing CID: %w", err)
		}
		str, err := multibase.Encode(multibase.Base58BTC, c.Hash())
		if err != nil {
			return fmt.Errorf("encoding multihash: %w", err)
		}
		fmt.Fprintln(cmd.OutOrStdout(), str)
		return nil
	},
}

func init() {
	Cmd.AddCommand(b58mhCmd)
}
