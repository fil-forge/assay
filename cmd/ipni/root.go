package ipni

import (
	"github.com/fil-forge/assay/cmd/ipni/metadata"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "ipni",
	Short: "Tools for inspecting IPNI data",
}

func init() {
	Cmd.AddCommand(metadata.Cmd)
}
