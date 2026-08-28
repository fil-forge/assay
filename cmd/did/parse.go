package did

import (
	"errors"
	"fmt"

	"github.com/fil-forge/assay/internal/peerkey"
	"github.com/fil-forge/ucantone/did"
	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:          "parse <did>",
	Short:        "Parse a DID.",
	Long:         "Parse a DID and print additional information where available (the libp2p peer ID for an Ed25519 did:key).",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := did.Parse(args[0])
		if err != nil {
			return fmt.Errorf("parsing DID: %w", err)
		}
		w := cmd.OutOrStdout()
		fmt.Fprintf(w, "DID:\t%s\n", d)

		if d.Method() != "key" {
			return nil
		}
		id, err := peerkey.DIDToPeerID(d)
		if err != nil {
			if errors.Is(err, peerkey.ErrUnsupportedKey) {
				return nil
			}
			return err
		}
		fmt.Fprintf(w, "PeerID:\t%s\n", id)
		return nil
	},
}

func init() {
	Cmd.AddCommand(parseCmd)
}
