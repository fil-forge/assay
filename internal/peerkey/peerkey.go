// Package peerkey converts between did:key DIDs and libp2p peer IDs. Only
// Ed25519 keys are supported; other key types return ErrUnsupportedKey.
package peerkey

import (
	"errors"
	"fmt"

	"github.com/fil-forge/ucantone/did"
	"github.com/fil-forge/ucantone/multikey"
	ed25519verifier "github.com/fil-forge/ucantone/multikey/ed25519/verifier"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multicodec"
)

// ErrUnsupportedKey is returned for DIDs and peer IDs whose key type has no
// conversion.
var ErrUnsupportedKey = errors.New("unsupported key type")

// DIDToPeerID derives the libp2p peer ID for an Ed25519 did:key.
func DIDToPeerID(d did.DID) (peer.ID, error) {
	if d.Method() != "key" {
		return "", fmt.Errorf("%w: DID method %q", ErrUnsupportedKey, d.Method())
	}
	verifier, err := multikey.Parse(d.Identifier())
	if err != nil {
		return "", fmt.Errorf("parsing did:key identifier: %w", err)
	}
	if verifier.Code() != multicodec.Ed25519Pub {
		return "", fmt.Errorf("%w: %s", ErrUnsupportedKey, verifier.Code())
	}
	pub, err := crypto.UnmarshalEd25519PublicKey(verifier.Raw())
	if err != nil {
		return "", fmt.Errorf("unmarshaling ed25519 public key: %w", err)
	}
	return peer.IDFromPublicKey(pub)
}

// PeerIDToDID derives the did:key DID for a peer ID that embeds an Ed25519
// public key.
func PeerIDToDID(id peer.ID) (did.DID, error) {
	pub, err := id.ExtractPublicKey()
	if err != nil {
		return did.Undef, fmt.Errorf("extracting public key: %w", err)
	}
	if _, ok := pub.(*crypto.Ed25519PublicKey); !ok {
		return did.Undef, fmt.Errorf("%w: %s", ErrUnsupportedKey, pub.Type())
	}
	raw, err := pub.Raw()
	if err != nil {
		return did.Undef, fmt.Errorf("reading raw public key: %w", err)
	}
	verifier, err := ed25519verifier.FromRaw(raw)
	if err != nil {
		return did.Undef, fmt.Errorf("decoding ed25519 verifier: %w", err)
	}
	return verifier.KeyDID(), nil
}
