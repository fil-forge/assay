package metadata_test

import (
	"bytes"
	"encoding/base64"
	"net/url"
	"testing"
	"time"

	cmdmetadata "github.com/fil-forge/assay/cmd/ipni/metadata"
	"github.com/fil-forge/go-ipni-tools/pkg/metadata"
	"github.com/fil-forge/libforge/commands"
	"github.com/fil-forge/libforge/commands/assert"
	"github.com/fil-forge/libforge/digestutil"
	"github.com/fil-forge/ucantone/testutil"
	"github.com/fil-forge/ucantone/ucan"
	"github.com/fil-forge/ucantone/ucan/invocation"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
)

func exec(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	cmdmetadata.Cmd.SetOut(&stdout)
	cmdmetadata.Cmd.SetErr(&stderr)
	cmdmetadata.Cmd.SetArgs(args)
	err := cmdmetadata.Cmd.Execute()
	return stdout.String(), err
}

// Example from the original debugger README: an index claim whose claim CID is
// a regular (non-identity) CID.
func TestParseIndexClaim(t *testing.T) {
	out, err := exec(t, "parse", "gID4AaNhY9gqWCUAAXESID3BhgVqQvFHsTo4TBfqcMJiTfE7ttRC+ofLgqAW5IeaYWUAYWnYKlgmAAGCBBIgUX8wwCPTHNNAl3QqAdrM+j5wl47D87QWT+Ps3ZlKQ4Y=")
	require.NoError(t, err)
	require.Contains(t, out, "ID:\t0x3e0000 (index claim)\n")
	require.Contains(t, out, "Claim:\tbafyreib5ygdak2sc6fd3coryjql6u4gcmjg7co5w2rbpvb6lqkqbnzehti\n")
	require.Contains(t, out, "Index:\tbagbaierakf7tbqbd2mongqexoqvadwwm7i7hbf4oypz3ifsp4pwn3gkkioda\n")
	require.Contains(t, out, "Expiration:\tnone\n")
}

type locationFixture struct {
	encoded string
	args    assert.LocationArguments
	exp     time.Time
}

func locationCommitment(t *testing.T) locationFixture {
	issuer := testutil.RandomIssuer(t)
	space := testutil.RandomDID(t)
	content := testutil.RandomDigest(t)
	loc, err := url.Parse("https://node.forge.example/blob/" + digestutil.Format(content))
	require.NoError(t, err)
	exp := time.Now().Add(time.Hour).Truncate(time.Second)
	end := uint64(4095)
	args := assert.LocationArguments{
		Space:    space,
		Content:  content,
		Location: []commands.CborURL{commands.CborURL(*loc)},
		Range:    &assert.Range{Start: 0, End: &end},
	}

	inv, err := assert.Location.Invoke(issuer, space, &args, invocation.WithExpiration(ucan.UnixTimestamp(exp.Unix())))
	require.NoError(t, err)
	invBytes, err := invocation.Encode(inv)
	require.NoError(t, err)

	idmh, err := multihash.Sum(invBytes, multihash.IDENTITY, -1)
	require.NoError(t, err)
	claim := cid.NewCidV1(cid.DagCBOR, idmh)

	md := metadata.MetadataContext.New(&metadata.LocationCommitmentMetadata{
		Claim:      claim,
		Expiration: exp.Unix(),
	})
	mdBytes, err := md.MarshalBinary()
	require.NoError(t, err)
	return locationFixture{base64.StdEncoding.EncodeToString(mdBytes), args, exp}
}

func TestParseEmbeddedLocationCommitment(t *testing.T) {
	f := locationCommitment(t)
	out, err := exec(t, "parse", f.encoded)
	require.NoError(t, err)
	require.Contains(t, out, "ID:\t0x3e0002 (location commitment)\n")
	require.Contains(t, out, "  Command: /assert/location\n")
	require.Contains(t, out, "  Subject: "+f.args.Space.String()+"\n")
	require.Contains(t, out, "    Content: "+digestutil.Format(f.args.Content)+"\n")
	require.Contains(t, out, "      "+f.args.Location[0].URL().String()+"\n")
	require.Contains(t, out, "    Range: 0-4095\n")
	require.Contains(t, out, "  Expiration: "+f.exp.Local().String()+"\n")
	require.Contains(t, out, "Expiration:\t"+f.exp.Local().String()+"\n")
}

func TestParseEmbeddedLocationCommitmentJSON(t *testing.T) {
	f := locationCommitment(t)
	out, err := exec(t, "parse", "--json", f.encoded)
	require.NoError(t, err)
	require.Contains(t, out, `"cmd":"/assert/location"`)
	require.NotContains(t, out, "  Command:")
}

func TestParseInvalid(t *testing.T) {
	_, err := exec(t, "parse", "!!!")
	require.Error(t, err)
	_, err = exec(t, "parse", base64.StdEncoding.EncodeToString([]byte("garbage")))
	require.Error(t, err)
}
