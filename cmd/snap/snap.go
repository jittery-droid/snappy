package snap

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagOptionDeployment = "deployment"
	flagOptionSeconds    = "duration"
	flagOptionProfile    = "profile"

	statusRunning = "Running"
)

func Cmds() []*cobra.Command {
	return []*cobra.Command{
		snapCmd,
	}
}

var snapCmd = &cobra.Command{
	Use:   "snap",
	Short: "get a snapshot for a pod in a given deployment",
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := exec.Command("bash", "-c", fmt.Sprintf("kubectl get pods -A | grep %s", flagDeployment)).Output()
		if err != nil {
			return err
		}
		podsList := strings.Split(string(out), "\n")
		var podName string
		for _, p := range podsList {
			fields := strings.Fields(p)
			fmt.Println(fields)
			if fields[3] == statusRunning {
				podName = fields[1]
				break
			}
		}

		pfCmd := exec.Command("bash", "-c",
			fmt.Sprintf("kubectl port-forward pod/%s 6060:6060 --namespace=roadie", podName))

		errChan := make(chan error, 1)
		forwardChan := make(chan struct{}, 1)

		go func() {
			err := pfCmd.Start()
			if err != nil {
				errChan <- err
			}
			fmt.Printf("port forwarding to pod %s on port 6060\n", podName)
			forwardChan <- struct{}{}
		}()

		select {
		case err = <-errChan:
			fmt.Println("start err")
			return err
		case <-forwardChan:
			break
		}

		var cmdStr, fileName string
		switch flagProfile {
		case "heap":
			fileName = "heap.out"
			cmdStr = "curl http://localhost:6060/debug/pprof/heap > " + fileName
		case "profile":
			fileName = "profile.out"
			cmdStr = fmt.Sprintf("curl http://localhost:6060/debug/pprof/profile?seconds=%d > ", flagSeconds) + fileName
		case "block":
			fileName = "block.out"
			cmdStr = "curl http://localhost:6060/debug/pprof/block > " + fileName
		case "mutex":
			fileName = "mutex.out"
			cmdStr = "curl http://localhost:6060/debug/pprof/mutex > " + fileName
		case "trace":
			fileName = "trace.out"
			cmdStr = fmt.Sprintf("curl http://localhost:6060/debug/pprof/trace?seconds=%d > ", flagSeconds) + fileName
		default:
			fmt.Println("invalid profile choice")
			return nil
		}

		_, err = exec.Command("bash", "-c", cmdStr).Output()
		if err != nil {
			fmt.Println("pprof err")
			return err
		}

		fmt.Printf("captured %s profile in %s", flagProfile, fileName)

		err = pfCmd.Process.Kill()
		if err != nil {
			fmt.Println("kill err")
			return err
		}

		return nil
	},
}

var (
	flagDeployment string
	flagProfile    string
	flagSeconds    int
)

func init() {
	snapCmd.Flags().StringVarP(&flagDeployment, flagOptionDeployment, "d", "", "search for pods within a deployment and forward to one")
	snapCmd.MarkFlagRequired(flagOptionDeployment)
	viper.BindPFlag(flagOptionDeployment, snapCmd.Flags().Lookup(flagOptionDeployment))

	snapCmd.Flags().StringVarP(&flagProfile, flagOptionProfile, "p", "", "choose a profile to capture from a pod in the deployment")
	snapCmd.MarkFlagRequired(flagOptionProfile)
	viper.BindPFlag(flagOptionProfile, snapCmd.Flags().Lookup(flagOptionProfile))

	snapCmd.Flags().IntVarP(&flagSeconds, flagOptionSeconds, "s", 30, "duration in which to capture a profile")
	viper.BindPFlag(flagOptionSeconds, snapCmd.Flags().Lookup(flagOptionSeconds))
}
