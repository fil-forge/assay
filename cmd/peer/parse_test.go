package peer_test

import (
	"bytes"
	"testing"

	cmdpeer "github.com/fil-forge/assay/cmd/peer"
	"github.com/stretchr/testify/require"
)

func exec(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	cmdpeer.Cmd.SetOut(&stdout)
	cmdpeer.Cmd.SetErr(&stderr)
	cmdpeer.Cmd.SetArgs(args)
	err := cmdpeer.Cmd.Execute()
	return stdout.String(), err
}

func TestParseEd25519(t *testing.T) {
	out, err := exec(t, "parse", "12D3KooWP11ydH5MT96hLTERNpyn2gPGJAkz5dH6PpNtFMgBRqXn")
	require.NoError(t, err)
	require.Equal(t, "PeerID:\t12D3KooWP11ydH5MT96hLTERNpyn2gPGJAkz5dH6PpNtFMgBRqXn\nDID:\tdid:key:z6MksdurUPk5bnKg34ZhRoCH6yrRdQzCTpSy8EvnheVJT7bE\n", out)
}

func TestParseRSA(t *testing.T) {
	out, err := exec(t, "parse", "QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N")
	require.NoError(t, err)
	require.Equal(t, "PeerID:\tQmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N\n", out)
}

func TestParseInvalid(t *testing.T) {
	_, err := exec(t, "parse", "nope")
	require.Error(t, err)
}
