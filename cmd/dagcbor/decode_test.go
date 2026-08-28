package dagcbor_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/fil-forge/assay/cmd/dagcbor"
	"github.com/fil-forge/ucantone/ipld/datamodel"
	"github.com/stretchr/testify/require"
)

func exec(t *testing.T, stdin []byte, args ...string) (string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	dagcbor.Cmd.SetOut(&stdout)
	dagcbor.Cmd.SetErr(&stderr)
	dagcbor.Cmd.SetIn(bytes.NewReader(stdin))
	dagcbor.Cmd.SetArgs(args)
	err := dagcbor.Cmd.Execute()
	return stdout.String(), err
}

func encode(t *testing.T) []byte {
	var buf bytes.Buffer
	m := datamodel.Map{"name": "assay", "count": int64(3), "tags": []any{"a", "b"}}
	require.NoError(t, m.MarshalCBOR(&buf))
	return buf.Bytes()
}

func TestDecodeFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.cbor")
	require.NoError(t, os.WriteFile(path, encode(t), 0644))

	out, err := exec(t, nil, "decode", path)
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal([]byte(out), &decoded), "output is plain JSON when stdout is not a terminal")
	require.Equal(t, "assay", decoded["name"])
	require.Equal(t, float64(3), decoded["count"])
}

func TestDecodeStdin(t *testing.T) {
	out, err := exec(t, encode(t), "decode")
	require.NoError(t, err)
	require.Contains(t, out, `"name": "assay"`)
}

func TestDecodeInvalid(t *testing.T) {
	_, err := exec(t, []byte{0xff, 0xff}, "decode")
	require.Error(t, err)
}
