package metadata

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/fil-forge/go-ipni-tools/pkg/metadata"
	"github.com/fil-forge/libforge/commands/assert"
	"github.com/fil-forge/libforge/digestutil"
	"github.com/fil-forge/ucantone/ucan"
	"github.com/fil-forge/ucantone/ucan/invocation"
	"github.com/ipfs/go-cid"
	ipnimd "github.com/ipni/go-libipni/metadata"
	"github.com/multiformats/go-multicodec"
	"github.com/multiformats/go-multihash"
	"github.com/spf13/cobra"
)

var parseJSON bool

var parseCmd = &cobra.Command{
	Use:   "parse <base64-metadata>",
	Short: "Parse IPNI metadata.",
	Long: `Parse base64 encoded IPNI metadata, as found in an IPNI provider result.
Content claim protocols are decoded; when the claim CID is an identity CID the
embedded claim invocation is decoded and printed too.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		metadataBytes, err := base64.StdEncoding.DecodeString(args[0])
		if err != nil {
			return fmt.Errorf("decoding base64: %w", err)
		}
		md := metadata.MetadataContext.New()
		if err := md.UnmarshalBinary(metadataBytes); err != nil {
			return fmt.Errorf("decoding metadata: %w", err)
		}
		return printMetadata(cmd.OutOrStdout(), md, parseJSON)
	},
}

func init() {
	parseCmd.Flags().BoolVarP(&parseJSON, "json", "j", false, "Print an embedded claim invocation as DAG-JSON instead of decoded arguments.")
	Cmd.AddCommand(parseCmd)
}

func printMetadata(w io.Writer, md ipnimd.Metadata, asJSON bool) error {
	// the provider may list one or more protocols for this CID; in the Forge
	// network the protocols are different kinds of content claim.
	for _, code := range md.Protocols() {
		protocol := md.Get(code)
		fmt.Fprintf(w, "ID:\t%#x (%s)\n", uint64(protocol.ID()), protocolName(protocol.ID()))

		hasClaim, ok := protocol.(metadata.HasClaim)
		if !ok {
			if protocolName(protocol.ID()) == "unknown" {
				fmt.Fprintln(w, "UNKNOWN METADATA")
			}
			continue
		}
		claim := hasClaim.GetClaim()
		fmt.Fprintf(w, "Claim:\t%s\n", claim)

		if claim.Prefix().MhType == multihash.IDENTITY {
			if err := printEmbeddedClaim(w, claim, asJSON); err != nil {
				return err
			}
		}

		switch typed := protocol.(type) {
		case *metadata.EqualsClaimMetadata:
			fmt.Fprintf(w, "Equals:\t%s\n", typed.Equals)
			printExpiration(w, "", typed.Expiration)
		case *metadata.IndexClaimMetadata:
			fmt.Fprintf(w, "Index:\t%s\n", typed.Index)
			printExpiration(w, "", typed.Expiration)
		case *metadata.LocationCommitmentMetadata:
			if typed.Shard != nil {
				fmt.Fprintf(w, "Shard:\t%s\n", typed.Shard)
			}
			if typed.Range != nil {
				if typed.Range.Length != nil {
					fmt.Fprintf(w, "Range:\t%d-%d (%d bytes)\n", typed.Range.Offset, typed.Range.Offset+*typed.Range.Length, *typed.Range.Length)
				} else {
					fmt.Fprintf(w, "Range:\t%d-\n", typed.Range.Offset)
				}
			}
			printExpiration(w, "", typed.Expiration)
		}
	}
	return nil
}

func protocolName(id multicodec.Code) string {
	switch id {
	case multicodec.TransportBitswap, multicodec.TransportIpfsGatewayHttp:
		return id.String()
	case metadata.LocationCommitmentID:
		return "location commitment"
	case metadata.EqualsClaimID:
		return "equivalency claim"
	case metadata.IndexClaimID:
		return "index claim"
	default:
		return "unknown"
	}
}

// printEmbeddedClaim decodes the claim invocation carried inline in an
// identity CID and prints it indented under the Claim line.
func printEmbeddedClaim(w io.Writer, claim cid.Cid, asJSON bool) error {
	dmh, err := multihash.Decode(claim.Hash())
	if err != nil {
		return fmt.Errorf("decoding claim multihash: %w", err)
	}
	inv, err := invocation.Decode(dmh.Digest)
	if err != nil {
		return fmt.Errorf("decoding embedded claim invocation: %w", err)
	}

	if asJSON {
		var buf bytes.Buffer
		if err := inv.MarshalDagJSON(&buf); err != nil {
			return fmt.Errorf("encoding claim as DAG-JSON: %w", err)
		}
		fmt.Fprintf(w, "  %s\n", buf.String())
		return nil
	}

	fmt.Fprintf(w, "  Command: %s\n", inv.Command())
	fmt.Fprintf(w, "  Subject: %s\n", inv.Subject())
	fmt.Fprintf(w, "  Issuer: %s\n", inv.Issuer())
	if err := printArguments(w, inv); err != nil {
		return err
	}
	if inv.Expiration() != nil {
		fmt.Fprintf(w, "  Expiration: %s\n", time.Unix(int64(*inv.Expiration()), 0))
	}
	return nil
}

// printArguments decodes the invocation arguments for the claim commands the
// Forge network publishes to IPNI. Other commands print nothing extra.
func printArguments(w io.Writer, inv ucan.Invocation) error {
	switch inv.Command().String() {
	case assert.Location.Command.String():
		var args assert.LocationArguments
		if err := args.UnmarshalCBOR(bytes.NewReader(inv.ArgumentsBytes())); err != nil {
			return fmt.Errorf("decoding location arguments: %w", err)
		}
		fmt.Fprintln(w, "  Arguments:")
		fmt.Fprintf(w, "    Space: %s\n", args.Space)
		fmt.Fprintf(w, "    Content: %s\n", digestutil.Format(args.Content))
		fmt.Fprintln(w, "    Location:")
		for _, l := range args.Location {
			fmt.Fprintf(w, "      %s\n", l.URL())
		}
		if args.Range != nil {
			if args.Range.End != nil {
				fmt.Fprintf(w, "    Range: %d-%d\n", args.Range.Start, *args.Range.End)
			} else {
				fmt.Fprintf(w, "    Range: %d-\n", args.Range.Start)
			}
		}
	case assert.Index.Command.String():
		var args assert.IndexArguments
		if err := args.UnmarshalCBOR(bytes.NewReader(inv.ArgumentsBytes())); err != nil {
			return fmt.Errorf("decoding index arguments: %w", err)
		}
		fmt.Fprintln(w, "  Arguments:")
		fmt.Fprintf(w, "    Index: %s\n", args.Index)
	case assert.Equals.Command.String():
		var args assert.EqualsArguments
		if err := args.UnmarshalCBOR(bytes.NewReader(inv.ArgumentsBytes())); err != nil {
			return fmt.Errorf("decoding equals arguments: %w", err)
		}
		fmt.Fprintln(w, "  Arguments:")
		fmt.Fprintf(w, "    Content: %s\n", digestutil.Format(args.Content))
		fmt.Fprintf(w, "    Equals: %s\n", args.Equals)
	}
	return nil
}

func printExpiration(w io.Writer, indent string, exp int64) {
	if exp > 0 {
		fmt.Fprintf(w, "%sExpiration:\t%s\n", indent, time.Unix(exp, 0))
	} else {
		fmt.Fprintf(w, "%sExpiration:\tnone\n", indent)
	}
}
