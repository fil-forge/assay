package did

import (
	"errors"
	"fmt"

	"github.com/fil-forge/assay/internal/peerkey"
	"github.com/fil-forge/ucantone/did"
	"github.com/fil-forge/ucantone/did/web"
	"github.com/spf13/cobra"
)

var resolveInsecure bool

var resolveCmd = &cobra.Command{
	Use:          "resolve <did>",
	Short:        "Resolve a did:web and print its keys.",
	Long:         "Fetch the DID document of a did:web and print each key it lists as a did:key, with the libp2p peer ID for an Ed25519 key. A Forge storage node is a did:web, and this is how to find the peer ID it advertises under.",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		d, err := web.Parse(args[0])
		if err != nil {
			return fmt.Errorf("parsing DID: %w", err)
		}
		resolver, err := web.NewResolver(web.WithInsecure(resolveInsecure))
		if err != nil {
			return fmt.Errorf("creating resolver: %w", err)
		}
		doc, err := resolver.Resolve(cmd.Context(), d)
		if err != nil {
			return fmt.Errorf("resolving %s: %w", d, err)
		}

		w := cmd.OutOrStdout()
		fmt.Fprintf(w, "DID:\t%s\n", doc.ID)
		if doc.VerificationMethods == nil {
			return nil
		}
		for _, vm := range doc.VerificationMethods.All() {
			// Only Multikey methods have a did:key to show: it is the method's
			// public key with the did:key prefix. Any other type is passed over,
			// as is a Multikey whose material is not a key.
			if vm.Type != did.MultikeyVerificationMethodType {
				continue
			}
			multibase, _ := vm.Material[did.MultikeyPublicKeyMultibaseProp].(string)
			key, err := did.New("key", multibase)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "Key:\t%s\n", key)

			id, err := peerkey.DIDToPeerID(key)
			if err != nil {
				if errors.Is(err, peerkey.ErrUnsupportedKey) {
					continue
				}
				return err
			}
			fmt.Fprintf(w, "PeerID:\t%s\n", id)
		}
		return nil
	},
}

func init() {
	resolveCmd.Flags().BoolVar(&resolveInsecure, "insecure", false, "fetch the DID document over HTTP, for a local stack")
	Cmd.AddCommand(resolveCmd)
}
