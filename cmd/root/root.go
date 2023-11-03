package root

import (
	"github.com/jittery-droid/snappy/cmd/kube"

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
	rootCmd.AddCommand(kube.Cmds()...)
}
