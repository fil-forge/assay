package peerkey_test

import (
	"testing"

	"github.com/fil-forge/assay/internal/peerkey"
	"github.com/fil-forge/ucantone/did"
	"github.com/fil-forge/ucantone/multikey/ed25519"
	"github.com/fil-forge/ucantone/multikey/secp256k1"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/require"
)

// Known pair from the original debugger README.
const (
	knownDID    = "did:key:z6MksdurUPk5bnKg34ZhRoCH6yrRdQzCTpSy8EvnheVJT7bE"
	knownPeerID = "12D3KooWP11ydH5MT96hLTERNpyn2gPGJAkz5dH6PpNtFMgBRqXn"
)

func TestKnownPair(t *testing.T) {
	id, err := peerkey.DIDToPeerID(did.MustParse(knownDID))
	require.NoError(t, err)
	require.Equal(t, knownPeerID, id.String())

	d, err := peerkey.PeerIDToDID(id)
	require.NoError(t, err)
	require.Equal(t, knownDID, d.String())
}

func TestRoundTrip(t *testing.T) {
	signer, err := ed25519.Generate()
	require.NoError(t, err)

	id, err := peerkey.DIDToPeerID(signer.KeyDID())
	require.NoError(t, err)

	d, err := peerkey.PeerIDToDID(id)
	require.NoError(t, err)
	require.Equal(t, signer.KeyDID(), d)
}

func TestUnsupported(t *testing.T) {
	_, err := peerkey.DIDToPeerID(did.MustParse("did:web:forge.example"))
	require.ErrorIs(t, err, peerkey.ErrUnsupportedKey)

	signer, err := secp256k1.Generate()
	require.NoError(t, err)
	_, err = peerkey.DIDToPeerID(signer.KeyDID())
	require.ErrorIs(t, err, peerkey.ErrUnsupportedKey)

	// RSA peer IDs hash the key instead of embedding it.
	rsaID, err := peer.Decode("QmYyQSo1c1Ym7orWxLYvCrM2EmxFTANf8wXmmE7DWjhx5N")
	require.NoError(t, err)
	_, err = peerkey.PeerIDToDID(rsaID)
	require.Error(t, err)
}
