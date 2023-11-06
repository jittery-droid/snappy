package root

import (
	"github.com/jittery-droid/snappy/cmd/kube"
	"github.com/jittery-droid/snappy/cmd/snap"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "snappy",
	Short: "get pprof snapshots easier",
	Long:  "snappy is a wrapper around pprof requests that will make profiling easier",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cmds := make([]*cobra.Command, 0)
	cmds = append(cmds, snap.Cmds()...)
	cmds = append(cmds, kube.Cmds()...)

	rootCmd.AddCommand(cmds...)
}
