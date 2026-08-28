package cid_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	cmdcid "github.com/fil-forge/assay/cmd/cid"
	gocid "github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
)

func exec(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	cmdcid.Cmd.SetOut(&stdout)
	cmdcid.Cmd.SetErr(&stderr)
	cmdcid.Cmd.SetArgs(args)
	err := cmdcid.Cmd.Execute()
	return stdout.String(), err
}

const v1 = "baguqeera3nw74wxxypvzcrf7jlieskytzaohhd4zyana2ndnk4ubjl2p5rvq"

const v1Info = `CID:        baguqeera3nw74wxxypvzcrf7jlieskytzaohhd4zyana2ndnk4ubjl2p5rvq (base32)
Version:    1
IPLD Codec: 0x129 (dag-json)
Digest:     zQmd7DGrjv8juYkDPsc4HHfCH6QUaMuQr32aHb4AuJWBZQv (base58btc)
Code:       0x12 (sha2-256)
Length:     32
`

func TestParse(t *testing.T) {
	out, err := exec(t, "parse", v1)
	require.NoError(t, err)
	require.Equal(t, v1Info, out)
}

func TestParseV0(t *testing.T) {
	out, err := exec(t, "parse", "QmdfTbBqBPQ7VNxZEYEj14VmRuZBkqFbiwReogJgS1zR1n")
	require.NoError(t, err)
	require.Contains(t, out, "CIDv0:      QmdfTbBqBPQ7VNxZEYEj14VmRuZBkqFbiwReogJgS1zR1n (base58btc)\n")
	require.Contains(t, out, "CIDv1:      bafybeihdwdcefgh4dqkjv67uzcmw7ojee6xedzdetojuzjevtenxquvyku (base32)\n")
	require.Contains(t, out, "IPLD Codec: 0x70 (dag-pb)\n")
}

func TestParseInvalid(t *testing.T) {
	_, err := exec(t, "parse", "not-a-cid")
	require.Error(t, err)
}

func TestDecode(t *testing.T) {
	c := gocid.MustParse(v1)
	path := filepath.Join(t.TempDir(), "cid.bin")
	require.NoError(t, os.WriteFile(path, c.Bytes(), 0644))

	out, err := exec(t, "decode", path)
	require.NoError(t, err)
	require.Equal(t, v1Info, out)
}

func TestB58MH(t *testing.T) {
	out, err := exec(t, "b58mh", "bafybeihfcdsxbirwgcuyrvabs3xi5adpootdutwy4zmluarltyskmszvla")
	require.NoError(t, err)
	require.Equal(t, "zQmdkq1rExA72pPVr83RZNA1uRtNNimATafgTJxbf2fW54K\n", out)
}
