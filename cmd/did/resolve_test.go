package did_test

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// serveDocument starts a server that answers every request with a DID
// document listing the given keys, and returns its did:web. A did:web names a
// host, never an IP address, so the server is addressed as localhost.
func serveDocument(t *testing.T, keys ...string) string {
	t.Helper()
	var id string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods := make([]map[string]any, 0, len(keys))
		for i, k := range keys {
			methods = append(methods, map[string]any{
				"id":                 fmt.Sprintf("%s#key-%d", id, i),
				"type":               "Multikey",
				"controller":         id,
				"publicKeyMultibase": k,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"@context":           []string{"https://www.w3.org/ns/did/v1"},
			"id":                 id,
			"verificationMethod": methods,
		})
	}))
	t.Cleanup(srv.Close)
	_, port, err := net.SplitHostPort(srv.Listener.Addr().String())
	require.NoError(t, err)
	id = "did:web:localhost%3A" + port
	return id
}

func TestResolveWeb(t *testing.T) {
	id := serveDocument(t, "z6MksdurUPk5bnKg34ZhRoCH6yrRdQzCTpSy8EvnheVJT7bE")
	out, err := exec(t, "resolve", "--insecure", id)
	require.NoError(t, err)
	require.Equal(t, "DID:\t"+id+"\n"+
		"Key:\tdid:key:z6MksdurUPk5bnKg34ZhRoCH6yrRdQzCTpSy8EvnheVJT7bE\n"+
		"PeerID:\t12D3KooWP11ydH5MT96hLTERNpyn2gPGJAkz5dH6PpNtFMgBRqXn\n", out)
}

func TestResolveWebListsEveryKey(t *testing.T) {
	id := serveDocument(t,
		"z6MksdurUPk5bnKg34ZhRoCH6yrRdQzCTpSy8EvnheVJT7bE",
		"z6Mkod5Jr3yd5SC7UDueqK4dAAw5xYJYjksy722tA9Boxc4z",
	)
	out, err := exec(t, "resolve", "--insecure", id)
	require.NoError(t, err)
	require.Contains(t, out, "Key:\tdid:key:z6MksdurUPk5bnKg34ZhRoCH6yrRdQzCTpSy8EvnheVJT7bE\n")
	require.Contains(t, out, "Key:\tdid:key:z6Mkod5Jr3yd5SC7UDueqK4dAAw5xYJYjksy722tA9Boxc4z\n")
	require.Equal(t, 2, strings.Count(out, "PeerID:\t"))
}

func TestResolveRefusesHTTPUnlessInsecure(t *testing.T) {
	id := serveDocument(t, "z6MksdurUPk5bnKg34ZhRoCH6yrRdQzCTpSy8EvnheVJT7bE")
	_, err := exec(t, "resolve", "--insecure=false", id)
	require.Error(t, err)
}

func TestResolveRefusesOtherMethods(t *testing.T) {
	_, err := exec(t, "resolve", "did:key:z6MksdurUPk5bnKg34ZhRoCH6yrRdQzCTpSy8EvnheVJT7bE")
	require.Error(t, err)
}
