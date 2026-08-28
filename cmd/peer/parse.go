package peer

import (
	"errors"
	"fmt"

	"github.com/fil-forge/assay/internal/peerkey"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"
)

var parseCmd = &cobra.Command{
	Use:          "parse <peer-id>",
	Short:        "Parse a libp2p peer ID.",
	Long:         "Parse a libp2p peer ID and print additional information where available (the did:key for an embedded Ed25519 public key).",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := peer.Decode(args[0])
		if err != nil {
			return fmt.Errorf("parsing peer ID: %w", err)
		}
		w := cmd.OutOrStdout()
		fmt.Fprintf(w, "PeerID:\t%s\n", id)

		d, err := peerkey.PeerIDToDID(id)
		if err != nil {
			// peer IDs that hash the key rather than embed it, or that embed a
			// key type other than Ed25519, have no did:key to show.
			if errors.Is(err, peerkey.ErrUnsupportedKey) || errors.Is(err, peer.ErrNoPublicKey) {
				return nil
			}
			return err
		}
		fmt.Fprintf(w, "DID:\t%s\n", d)
		return nil
	},
}

func init() {
	Cmd.AddCommand(parseCmd)
}
