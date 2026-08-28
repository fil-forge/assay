package flatfs

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "flatfs",
	Short: "Tools for inspecting FlatFS blob stores",
}
