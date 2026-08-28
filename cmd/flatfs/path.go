package flatfs

import (
	"fmt"
	"path"
	"strings"

	"github.com/fil-forge/libforge/digestutil"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multihash"
	"github.com/spf13/cobra"
)

var pathCmd = &cobra.Command{
	Use:   "path <blob-cid-or-multihash>",
	Short: "Print the FlatFS path for a blob (Piri layout).",
	Long: `Given a blob CID or a multibase encoded multihash, print the path of the blob
within a Piri FlatFS blob store: a next-to-last/2 shard directory and a base32
key with a ".data" suffix.`,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		digest, err := parseDigest(args[0])
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), Path(digest))
		return nil
	},
}

func init() {
	Cmd.AddCommand(pathCmd)
}

// parseDigest accepts a CID (any version) or a multibase encoded multihash
// and returns the multihash. A CID is tried first: a bare sha2-256 multihash
// also parses as a CIDv0 and yields the same digest, so this order is safe.
func parseDigest(input string) (multihash.Multihash, error) {
	if c, err := cid.Decode(input); err == nil {
		return c.Hash(), nil
	}
	digest, err := digestutil.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as a CID or multihash: %w", input, err)
	}
	return digest, nil
}

// Path returns the store-relative path of a blob in a Piri FlatFS blob store.
// It mirrors piri's Base32FlatFSKeyEncoder: the base32 multibase encoding of
// the multihash with the base prefix stripped, sharded by NextToLast(2).
func Path(digest multihash.Multihash) string {
	b32, _ := multibase.Encode(multibase.Base32, digest)
	key := b32[1:] // strip the multibase prefix
	return path.Join("/", nextToLast(key, 2), key+".data")
}

// nextToLast is the FlatFS "next-to-last" sharding function: the suffixLen
// characters before the last character of the key, padded with underscores
// when the key is too short.
func nextToLast(key string, suffixLen int) string {
	str := strings.Repeat("_", suffixLen+1) + key
	offset := len(str) - suffixLen - 1
	return str[offset : offset+suffixLen]
}
