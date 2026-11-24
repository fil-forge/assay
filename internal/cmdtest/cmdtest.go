// Package cmdtest holds helpers for tests that run a command.
package cmdtest

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ResetFlags puts every flag of cmd and its subcommands back to its default
// when the test ends. A command keeps the flag values of its last run, so
// without this a flag one test sets is still set for whichever test runs next.
func ResetFlags(t *testing.T, cmd *cobra.Command) {
	t.Helper()
	t.Cleanup(func() { resetFlags(t, cmd) })
}

func resetFlags(t *testing.T, cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if err := f.Value.Set(f.DefValue); err != nil {
			t.Errorf("resetting flag --%s of %s: %v", f.Name, cmd.CommandPath(), err)
		}
		f.Changed = false
	})
	for _, sub := range cmd.Commands() {
		resetFlags(t, sub)
	}
}
