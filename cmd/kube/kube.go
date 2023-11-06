package kube

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const flagOptionService = "service"

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		lsCmd,
	}
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list deployments for a service",
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := exec.Command("bash", "-c", fmt.Sprintf("kubectl get deployments -A | grep %s", flagService)).Output()
		if err != nil {
			return err
		}

		fmt.Println(string(out))
		return nil
	},
}

var flagService string

func init() {
	lsCmd.Flags().StringVarP(&flagService, flagOptionService, "s", "", "search for deployments of this service")
	lsCmd.MarkFlagRequired(flagOptionService)
	viper.BindPFlag(flagOptionService, lsCmd.Flags().Lookup(flagOptionService))
}
