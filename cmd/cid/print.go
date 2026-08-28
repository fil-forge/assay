package cid

import (
	"fmt"
	"io"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multicodec"
	"github.com/multiformats/go-multihash"
)

func printCID(w io.Writer, c cid.Cid) error {
	mh, err := multihash.Decode(c.Hash())
	if err != nil {
		return fmt.Errorf("decoding multihash: %w", err)
	}
	digestStr, _ := multibase.Encode(multibase.Base58BTC, c.Hash())
	digestBaseStr := multibase.EncodingToStr[multibase.Base58BTC]

	if c.Prefix().Version == 0 {
		fmt.Fprintf(w, "%-12s%s (%s)\n", "CIDv0:", c.String(), multibase.EncodingToStr[multibase.Base58BTC])
		cv1 := cid.NewCidV1(uint64(multicodec.DagPb), c.Hash())
		cidv1Str, _ := cv1.StringOfBase(multibase.Base32)
		fmt.Fprintf(w, "%-12s%s (%s)\n", "CIDv1:", cidv1Str, multibase.EncodingToStr[multibase.Base32])
	} else {
		cidStr, _ := c.StringOfBase(multibase.Base32)
		fmt.Fprintf(w, "%-12s%s (%s)\n", "CID:", cidStr, multibase.EncodingToStr[multibase.Base32])
	}
	fmt.Fprintf(w, "%-12s%d\n", "Version:", c.Prefix().Version)
	fmt.Fprintf(w, "%-12s%#x (%s)\n", "IPLD Codec:", c.Prefix().Codec, multicodec.Code(c.Prefix().Codec).String())
	fmt.Fprintf(w, "%-12s%s (%s)\n", "Digest:", digestStr, digestBaseStr)
	fmt.Fprintf(w, "%-12s%#x (%s)\n", "Code:", mh.Code, multicodec.Code(mh.Code).String())
	fmt.Fprintf(w, "%-12s%d\n", "Length:", mh.Length)
	return nil
}
