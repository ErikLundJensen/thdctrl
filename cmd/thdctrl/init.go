package thdctrl

import (
	"fmt"
	"os"
	"strconv"
	"github.com/spf13/cobra"
	"github.com/eriklundjensen/thdctrl/pkg/robot"
	"github.com/eriklundjensen/thdctrl/pkg/hetznerapi"
)

var initCmdFlags struct {
		skipReboot bool
		enableRescueSystem bool
		disk string
		serverNumber int
}

var initCmd = &cobra.Command{
	Use:   "init <serverNumber>",
	Short: "Initialize the application",
	Args: cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
			serverNumber, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Printf("Error parsing server number: %v\n", err)
				return
			}
			initializeServer(RobotClient, initCmdFlags.skipReboot, initCmdFlags.enableRescueSystem, initCmdFlags.disk, serverNumber)
	},
}

func init() {
		initCmd.Flags().BoolVarP(&initCmdFlags.skipReboot, "skipReboot", "n", false ,"Skip reboot of server after enabling rescue system")
		initCmd.Flags().BoolVarP(&initCmdFlags.enableRescueSystem, "enable-rescue-system", "r", false, "enableRescueSystem: entering rescue system even if rescue system already enabled. This will generate a new password")
		initCmd.Flags().StringVarP(&initCmdFlags.disk, "disk", "d", "nvme0n1", "Disk to use for installation of image")
		addCommand(initCmd)
}

func initializeServer(client robot.Client, skipReboot bool, enableRescueSystem bool, disk string, serverNumber int) {
	sshPassword := os.Getenv("HETZNER_SSH_PASSWORD") // Set your Hetzner password in environment variable

	rescue, err := hetznerapi.GetRescueSystemDetails(client, serverNumber)
	if err != nil {
			fmt.Printf("Error getting rescue system status: %v\n", err)
			return
	}

	if (!rescue.Rescue.Active || enableRescueSystem) {
			rescue, err = hetznerapi.EnableRescueSystem(client, serverNumber)
			if err != nil {
					fmt.Printf("Error enabling rescue system: %v\n", err)
					return
			}
	}

	if !skipReboot  {
			err = hetznerapi.RebootServer(client, serverNumber)
	}
	if (err != nil || rescue == nil) {
			fmt.Printf("Rescue system state is not available: %v\n", err)
			return
	}
	sshClient := &hetznerapi.SSHClient {
			Host: rescue.Rescue.ServerIP,
			Port: "22",
	}
	sshUser := "root"
	if rescue.Rescue.Password != "" {
			sshPassword = rescue.Rescue.Password
	}
	sshClient.Auth(sshUser, sshPassword)

	sshClient.WaitForReboot()
	fmt.Printf("Server rebooted with Talos\n")

	// Don't use image hcloud-amd64 target Hetzner Cloud, use Talos 'metal' image instead
	version := "v1.9.2"
	imageUrl := fmt.Sprintf("https://github.com/siderolabs/talos/releases/download/%s/metal-amd64.raw.zst", version)
	output, err := sshClient.DownloadImage(imageUrl)
	if err != nil {
		fmt.Printf("Failed to download image: %v, output %s\n", err, output)
		return
	}

	output, err = sshClient.InstallImage(disk)
	if err != nil {
		fmt.Printf("Failed to install image: %v output %s\n", err, output)
		return
	}

	hetznerapi.RebootServer(client, serverNumber)

	// Wait for Talos API to become available
	// Apply Talos configuration
	// Apply Cilium
	// Reboot node
}