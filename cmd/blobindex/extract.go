package blobindex

import (
	"bytes"
	"fmt"
	"io"
	"slices"

	"github.com/fil-forge/assay/internal/cmdutil"
	"github.com/fil-forge/automobile"
	"github.com/fil-forge/libforge/blobindex"
	"github.com/fil-forge/libforge/digestutil"
	mh "github.com/multiformats/go-multihash"
	"github.com/spf13/cobra"
)

var extractJSON bool

var extractCmd = &cobra.Command{
	Use:          "extract [car-file]",
	Short:        "Extract a sharded DAG index from a CAR.",
	Long:         "Extract a sharded DAG index that has been archived to a CAR. Reads the file given, or stdin when none is.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		archive, err := cmdutil.ReadInput(cmd, args)
		if err != nil {
			return err
		}
		if extractJSON {
			return printBlocks(cmd.OutOrStdout(), archive)
		}
		index, err := blobindex.Extract(bytes.NewReader(archive))
		if err != nil {
			return fmt.Errorf("extracting blob index: %w", err)
		}
		printIndex(cmd.OutOrStdout(), index)
		return nil
	},
}

func init() {
	extractCmd.Flags().BoolVarP(&extractJSON, "json", "j", false, "Print every block in the CAR as DAG-JSON instead of the decoded index.")
	Cmd.AddCommand(extractCmd)
}

// printBlocks prints each block of the CAR as its CID followed by its
// DAG-JSON rendering.
func printBlocks(w io.Writer, archive []byte) error {
	_, blocks, err := automobile.Decode(bytes.NewReader(archive))
	if err != nil {
		return fmt.Errorf("decoding CAR: %w", err)
	}
	for _, b := range blocks {
		s, err := cmdutil.FormatCBOR(b.Data)
		if err != nil {
			return fmt.Errorf("formatting block %s: %w", b.Link, err)
		}
		fmt.Fprintf(w, "%s\n%s\n\n", b.Link, s)
	}
	return nil
}

// printIndex prints shards and their slices, sorted by digest so output is
// stable between runs.
func printIndex(w io.Writer, index blobindex.ShardedDagIndex) {
	shards := index.Shards()
	fmt.Fprintf(w, "Shards (%d):\n", shards.Size())
	for _, shard := range sortedDigests(shards.Keys()) {
		slices := shards.Get(shard)
		fmt.Fprintf(w, "  %s\n", digestutil.Format(shard))
		fmt.Fprintf(w, "    Slices (%d):\n", slices.Size())
		for _, digest := range sortedDigests(slices.Keys()) {
			rng := slices.Get(digest)
			fmt.Fprintf(w, "      %s @ %d-%d\n", digestutil.Format(digest), rng.Start, rng.End)
		}
	}
}

func sortedDigests(seq func(yield func(mh.Multihash) bool)) []mh.Multihash {
	return slices.SortedFunc(seq, func(a, b mh.Multihash) int { return bytes.Compare(a, b) })
}
