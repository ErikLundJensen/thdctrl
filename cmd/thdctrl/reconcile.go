package thdctrl

import (
	"fmt"
	"os"

	v1alpha1 "github.com/eriklundjensen/thdctrl/pkg/api/server/v1alpha"
	"github.com/eriklundjensen/thdctrl/pkg/controller"
  "github.com/eriklundjensen/thdctrl/pkg/hetznerapi"
	"github.com/eriklundjensen/thdctrl/pkg/robot"
	"github.com/spf13/cobra"
	yaml "github.com/goccy/go-yaml"
)

var (
	filename string

	reconcileCmd = &cobra.Command{
		Use:   "reconcile",
		Short: "Reconcile server configuration from file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if filename == "" {
				return fmt.Errorf("filename is required")
			}
      sshClient := &hetznerapi.SSHClient{}
			return reconcileFromFile(RobotClient, sshClient, filename)
		},
	}
)

func init() {
	reconcileCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename containing server configuration (required)")
	reconcileCmd.MarkFlagRequired("filename")
	addCommand(reconcileCmd)
}

func readServerConfig(filename string) (*v1alpha1.ServerParameters, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var server v1alpha1.ServerParameters
	if err := yaml.Unmarshal(data, &server); err != nil {
		return nil, fmt.Errorf("error parsing yaml: %v", err)
	}

	return &server, nil
}

func reconcileFromFile(client robot.ClientInterface, sshClient *hetznerapi.SSHClient, filename string) error {
	server, err := readServerConfig(filename)
	if err != nil {
		return err
	}

	fmt.Printf("Read configuration for server %d\n", server.ServerNumber)

	// Check server status
	status := controller.DetermineServerStatus(client, sshClient, server)
	fmt.Printf("Server status: %s\n", status)

	// Continue only if Talos API is available
	if status != controller.TalosAPIAvailable {
		return fmt.Errorf("server is not in expected state: %s", status)
	}

	return nil
}
