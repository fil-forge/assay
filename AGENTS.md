# assay

Inspection tools for Forge network data: CIDs, DIDs, libp2p peer IDs, DAG-CBOR,
blob indexes, Piri FlatFS paths and IPNI metadata. A cobra CLI, one command
group per directory under `cmd/`.

## Layout

| Path | What it is |
|---|---|
| `main.go` | calls `cmd.Execute()` |
| `cmd/root.go` | root command; `init()` registers every group's `Cmd` |
| `cmd/<group>/root.go` | exported `var Cmd = &cobra.Command{Use: "<group>"}` |
| `cmd/<group>/<leaf>.go` | one leaf command per file, self-registered in `init()`, test alongside |
| `internal/cmdutil` | `ReadInput` (file arg or stdin), `FormatCBOR` (DAG-CBOR to DAG-JSON) |
| `internal/peerkey` | Ed25519 `did:key` to libp2p peer ID and back |

## Conventions

- Leaf commands use `RunE` and `SilenceUsage: true`; return wrapped errors
  (`fmt.Errorf("...: %w", err)`), never `panic` or `cobra.CheckErr`.
- All output goes to `cmd.OutOrStdout()` so commands are testable in-process
  (set `Cmd.SetOut`/`SetIn`/`SetArgs`, call `Cmd.Execute()`; see any `*_test.go`).
- Flags are package globals bound with `BoolVarP` and friends.
- Pretty printing of IPLD values reuses `github.com/fil-forge/ucantool/pkg/ipldfmt`.
  Domain types come from `libforge` (blob index, `digestutil`, `commands/assert`),
  `ucantone` (DIDs, multikey, invocations) and `go-ipni-tools` (claim metadata).
  Do not add `go-ipld-prime` based code; the Go direction is cbor-gen.
- `ucantone`, `libforge`, `ucantool` and `go-ipni-tools` are pinned by
  pseudo-version; bump with `GOWORK=off go get <module>@<commit>` then
  `GOWORK=off go mod tidy`. The fil-forge `go.work` lists this module, so local
  builds resolve sibling clones.

## Build and test

```sh
go build ./... && go vet ./... && go test ./...
GOWORK=off go build ./...   # proves go.mod pins resolve without the workspace
```
