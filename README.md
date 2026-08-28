# assay

Inspection tools for Forge network data. An assay is the test that tells you
what a piece of metal is made of; this tool does the same for CIDs, DIDs, peer
IDs, DAG-CBOR blocks, blob indexes, Piri FlatFS paths and IPNI metadata.

## Install

```sh
go install github.com/fil-forge/assay@latest
```

## Usage

### `assay blobindex extract [car-file]`

Extract a sharded DAG index that has been archived to a CAR. Reads the file
given, or stdin when none is, so you can pipe into it. Slices are printed as
`digest @ start-end` (inclusive byte offsets within the shard).

```console
$ assay blobindex extract index.car
Shards (14):
  zQmeDoFFnXb2zYYi4hFdNKz9kAwu5eHYLm7LokxQyfmjp1M
    Slices (253):
      zQmeoqurjukfWtjsNry2ogGNWifv66RsKRCNGar8xvf5apS @ 12583437-13632013
      zQmW6zw35SdyRhbuV44Kjwwju54drMXHup8JjRMEzDYb8q8 @ 177215992-178264568
...
```

`--json` (`-j`) prints every block in the CAR as DAG-JSON instead.

### `assay cid b58mh <cid>`

Extract the multihash from a CID and print it as a multibase base58btc string.

```console
$ assay cid b58mh bafybeihfcdsxbirwgcuyrvabs3xi5adpootdutwy4zmluarltyskmszvla
zQmdkq1rExA72pPVr83RZNA1uRtNNimATafgTJxbf2fW54K
```

### `assay cid decode <file>`

Decode a byte encoded CID and print information about it.

```console
$ assay cid decode ./cid.bin
CID:        baguqeera3nw74wxxypvzcrf7jlieskytzaohhd4zyana2ndnk4ubjl2p5rvq (base32)
Version:    1
IPLD Codec: 0x129 (dag-json)
Digest:     zQmd7DGrjv8juYkDPsc4HHfCH6QUaMuQr32aHb4AuJWBZQv (base58btc)
Code:       0x12 (sha2-256)
Length:     32
```

### `assay cid parse <cid>`

Parse a string encoded CID and print information about it.

```console
$ assay cid parse baguqeera3nw74wxxypvzcrf7jlieskytzaohhd4zyana2ndnk4ubjl2p5rvq
CID:        baguqeera3nw74wxxypvzcrf7jlieskytzaohhd4zyana2ndnk4ubjl2p5rvq (base32)
Version:    1
IPLD Codec: 0x129 (dag-json)
Digest:     zQmd7DGrjv8juYkDPsc4HHfCH6QUaMuQr32aHb4AuJWBZQv (base58btc)
Code:       0x12 (sha2-256)
Length:     32
```

### `assay dagcbor decode [file]`

Decode DAG-CBOR data from a file or stdin, format it as DAG-JSON and print it.
Output is syntax highlighted when stdout is a terminal.

### `assay did parse <did>`

Parse a DID. For an Ed25519 `did:key` the matching libp2p peer ID is printed
as well.

```console
$ assay did parse did:key:z6MksdurUPk5bnKg34ZhRoCH6yrRdQzCTpSy8EvnheVJT7bE
DID:    did:key:z6MksdurUPk5bnKg34ZhRoCH6yrRdQzCTpSy8EvnheVJT7bE
PeerID: 12D3KooWP11ydH5MT96hLTERNpyn2gPGJAkz5dH6PpNtFMgBRqXn
```

### `assay flatfs path <blob-cid-or-multihash>`

Given a blob CID or a multibase encoded multihash, print the path of the blob
within a Piri FlatFS blob store.

```console
$ assay flatfs path bagbaiera3f6ylgq5yqop4scfkgxibtbfzfgegm5gzrssdxv37zdkmxz2j4vq
/6k/ciqns7mftio4ihh6jbcvdluazqs4stcdgotmyzjb32574rvgl45e6ky.data

$ assay flatfs path zQmQE3fBFV8xQKeA34gspkmoWyyvJKKJYzybPGq6zLu4pTw
/od/ciqbybsb3rvglosb2m6pfdnweid6dp7mk4tnp6mimilduuwa4v4hodq.data
```

### `assay ipni metadata parse <value>`

Parse base64 encoded IPNI metadata, as found in an IPNI provider result. When
the claim CID is an identity CID, the embedded claim invocation is decoded and
its arguments printed. `--json` (`-j`) prints the embedded claim as DAG-JSON
instead.

```console
$ assay ipni metadata parse gID4AaNhY9gqWCUAAXESID3BhgVqQvFHsTo4TBfqcMJiTfE7ttRC+ofLgqAW5IeaYWUAYWnYKlgmAAGCBBIgUX8wwCPTHNNAl3QqAdrM+j5wl47D87QWT+Ps3ZlKQ4Y=
ID:     0x3e0000 (index claim)
Claim:  bafyreib5ygdak2sc6fd3coryjql6u4gcmjg7co5w2rbpvb6lqkqbnzehti
Index:  bagbaierakf7tbqbd2mongqexoqvadwwm7i7hbf4oypz3ifsp4pwn3gkkioda
Expiration:     none
```

### `assay peer parse <libp2p-peer-id>`

Parse a libp2p peer ID. When the ID embeds an Ed25519 public key the matching
`did:key` is printed as well.

```console
$ assay peer parse 12D3KooWP11ydH5MT96hLTERNpyn2gPGJAkz5dH6PpNtFMgBRqXn
PeerID: 12D3KooWP11ydH5MT96hLTERNpyn2gPGJAkz5dH6PpNtFMgBRqXn
DID:    did:key:z6MksdurUPk5bnKg34ZhRoCH6yrRdQzCTpSy8EvnheVJT7bE
```

## License

Dual-licensed under [MIT + Apache 2.0](LICENSE.md)
