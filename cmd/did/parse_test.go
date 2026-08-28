package did_test

import (
	"bytes"
	"testing"

	cmddid "github.com/fil-forge/assay/cmd/did"
	"github.com/stretchr/testify/require"
)

func exec(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	cmddid.Cmd.SetOut(&stdout)
	cmddid.Cmd.SetErr(&stderr)
	cmddid.Cmd.SetArgs(args)
	err := cmddid.Cmd.Execute()
	return stdout.String(), err
}

func TestParseKey(t *testing.T) {
	out, err := exec(t, "parse", "did:key:z6MksdurUPk5bnKg34ZhRoCH6yrRdQzCTpSy8EvnheVJT7bE")
	require.NoError(t, err)
	require.Equal(t, "DID:\tdid:key:z6MksdurUPk5bnKg34ZhRoCH6yrRdQzCTpSy8EvnheVJT7bE\nPeerID:\t12D3KooWP11ydH5MT96hLTERNpyn2gPGJAkz5dH6PpNtFMgBRqXn\n", out)
}

func TestParseWeb(t *testing.T) {
	out, err := exec(t, "parse", "did:web:forge.example")
	require.NoError(t, err)
	require.Equal(t, "DID:\tdid:web:forge.example\n", out)
}

func TestParseInvalid(t *testing.T) {
	_, err := exec(t, "parse", "nope")
	require.Error(t, err)
}
