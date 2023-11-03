package kube

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagOptionService    = "service"
	flagOptionDeployment = "deployment"

	statusRunning = "Running"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		lsCmd,
		fwdCmd,
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

var fwdCmd = &cobra.Command{
	Use:   "fwd",
	Short: "list deployments for a service",
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := exec.Command("bash", "-c", fmt.Sprintf("kubectl get pods -A | grep %s", flagDeployment)).Output()
		if err != nil {
			return err
		}
		podsList := strings.Split(string(out), "\n")
		var podName string
		for _, p := range podsList {
			fields := strings.Fields(p)
			if fields[3] == statusRunning {
				podName = fields[1]
				break
			}
		}

		fmt.Printf("port forwarding to pod %s on port 6060\n", podName)
		out, err = exec.Command("bash", "-c", fmt.Sprintf("kubectl port-forward pod/%s 3000:3000 --namespace=roadie", podName)).Output()
		if err != nil {
			return err
		}
		fmt.Println(out)
		return nil
	},
}

var (
	flagService    string
	flagDeployment string
)

func init() {
	lsCmd.Flags().StringVarP(&flagService, flagOptionService, "s", "", "search for deployments of this service")
	lsCmd.MarkFlagRequired(flagOptionService)
	viper.BindPFlag(flagOptionService, lsCmd.Flags().Lookup(flagOptionService))

	fwdCmd.Flags().StringVarP(&flagDeployment, flagOptionDeployment, "d", "", "search for pods within a deployment and forward to one")
	fwdCmd.MarkFlagRequired(flagOptionDeployment)
	viper.BindPFlag(flagOptionDeployment, fwdCmd.Flags().Lookup(flagOptionDeployment))
}
