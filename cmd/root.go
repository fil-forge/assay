package cmd

import (
	"os"

	"github.com/fil-forge/assay/cmd/blobindex"
	cmdcid "github.com/fil-forge/assay/cmd/cid"
	"github.com/fil-forge/assay/cmd/dagcbor"
	cmddid "github.com/fil-forge/assay/cmd/did"
	"github.com/fil-forge/assay/cmd/flatfs"
	"github.com/fil-forge/assay/cmd/ipni"
	cmdpeer "github.com/fil-forge/assay/cmd/peer"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "assay",
	Short: "Inspection tools for the Forge network",
	Long:  `assay is a collection of tools for inspecting and debugging Forge network data.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(blobindex.Cmd)
	rootCmd.AddCommand(cmdcid.Cmd)
	rootCmd.AddCommand(dagcbor.Cmd)
	rootCmd.AddCommand(cmddid.Cmd)
	rootCmd.AddCommand(flatfs.Cmd)
	rootCmd.AddCommand(ipni.Cmd)
	rootCmd.AddCommand(cmdpeer.Cmd)
}
