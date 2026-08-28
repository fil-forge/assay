package flatfs_test

import (
	"bytes"
	"testing"

	"github.com/fil-forge/assay/cmd/flatfs"
	"github.com/fil-forge/libforge/digestutil"
	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
)

func exec(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	flatfs.Cmd.SetOut(&stdout)
	flatfs.Cmd.SetErr(&stderr)
	flatfs.Cmd.SetArgs(args)
	err := flatfs.Cmd.Execute()
	return stdout.String(), err
}

// Examples from the original debugger README.
func TestPathFromCID(t *testing.T) {
	out, err := exec(t, "path", "bagbaiera3f6ylgq5yqop4scfkgxibtbfzfgegm5gzrssdxv37zdkmxz2j4vq")
	require.NoError(t, err)
	require.Equal(t, "/6k/ciqns7mftio4ihh6jbcvdluazqs4stcdgotmyzjb32574rvgl45e6ky.data\n", out)
}

func TestPathFromDigest(t *testing.T) {
	out, err := exec(t, "path", "zQmQE3fBFV8xQKeA34gspkmoWyyvJKKJYzybPGq6zLu4pTw")
	require.NoError(t, err)
	require.Equal(t, "/od/ciqbybsb3rvglosb2m6pfdnweid6dp7mk4tnp6mimilduuwa4v4hodq.data\n", out)
}

func TestPathCIDAndDigestAgree(t *testing.T) {
	const digest = "zQmQE3fBFV8xQKeA34gspkmoWyyvJKKJYzybPGq6zLu4pTw"
	mh, err := digestutil.Parse(digest)
	require.NoError(t, err)
	rawCID := cid.NewCidV1(cid.Raw, mh).String()

	fromCID, err := exec(t, "path", rawCID)
	require.NoError(t, err)
	fromDigest, err := exec(t, "path", digest)
	require.NoError(t, err)
	require.Equal(t, fromDigest, fromCID)
}

func TestPathInvalid(t *testing.T) {
	_, err := exec(t, "path", "nope")
	require.Error(t, err)
}
