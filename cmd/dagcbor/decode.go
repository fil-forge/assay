package dagcbor

import (
	"fmt"

	"github.com/fil-forge/assay/internal/cmdutil"
	"github.com/spf13/cobra"
)

var decodeCmd = &cobra.Command{
	Use:          "decode [path]",
	Short:        "Decode DAG-CBOR data and print it as DAG-JSON.",
	Long:         "Decode DAG-CBOR encoded data from a file or stdin, format it as DAG-JSON and print it.",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := cmdutil.ReadInput(cmd, args)
		if err != nil {
			return err
		}
		s, err := cmdutil.FormatCBOR(b)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), s)
		return nil
	},
}

func init() {
	Cmd.AddCommand(decodeCmd)
}
