package blobindex_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	cmdblobindex "github.com/fil-forge/assay/cmd/blobindex"
	"github.com/fil-forge/libforge/blobindex"
	"github.com/fil-forge/libforge/digestutil"
	"github.com/fil-forge/ucantone/testutil"
	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/require"
)

func exec(t *testing.T, stdin []byte, args ...string) (string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	cmdblobindex.Cmd.SetOut(&stdout)
	cmdblobindex.Cmd.SetErr(&stderr)
	cmdblobindex.Cmd.SetIn(bytes.NewReader(stdin))
	cmdblobindex.Cmd.SetArgs(args)
	err := cmdblobindex.Cmd.Execute()
	return stdout.String(), err
}

type fixture struct {
	car   []byte
	shard multihash.Multihash
	slice multihash.Multihash
}

func archive(t *testing.T) fixture {
	shard := testutil.RandomDigest(t)
	slice := testutil.RandomDigest(t)
	index := blobindex.NewShardedDagIndex(1)
	index.SetSlice(shard, slice, blobindex.Range{Start: 100, End: 199})
	var buf bytes.Buffer
	require.NoError(t, index.Archive(&buf))
	return fixture{buf.Bytes(), shard, slice}
}

func TestExtractFile(t *testing.T) {
	f := archive(t)
	path := filepath.Join(t.TempDir(), "index.car")
	require.NoError(t, os.WriteFile(path, f.car, 0644))

	out, err := exec(t, nil, "extract", path)
	require.NoError(t, err)
	expected := fmt.Sprintf("Shards (1):\n  %s\n    Slices (1):\n      %s @ 100-199\n", digestutil.Format(f.shard), digestutil.Format(f.slice))
	require.Equal(t, expected, out)
}

func TestExtractStdin(t *testing.T) {
	f := archive(t)
	out, err := exec(t, f.car, "extract")
	require.NoError(t, err)
	require.Contains(t, out, digestutil.Format(f.slice)+" @ 100-199")
}

func TestExtractJSON(t *testing.T) {
	f := archive(t)
	out, err := exec(t, f.car, "extract", "--json")
	require.NoError(t, err)
	// one shard block plus the root block
	require.Equal(t, 2, bytes.Count([]byte(out), []byte("\n\n")))
	require.Contains(t, out, `"index/sharded/dag@0.1"`)
}

func TestExtractInvalid(t *testing.T) {
	_, err := exec(t, []byte("not a car"), "extract")
	require.Error(t, err)
}
