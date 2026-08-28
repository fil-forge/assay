package peer

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "peer",
	Short: "Tools for inspecting libp2p peer IDs",
}
