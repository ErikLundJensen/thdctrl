package thdctrl

import (
	"fmt"
	"os"
	"strconv"
	"github.com/spf13/cobra"
	"github.com/eriklundjensen/thdctrl/pkg/robot"
	"github.com/eriklundjensen/thdctrl/pkg/hetznerapi"
)

type cmdFlags struct {
		skipReboot bool
		enableRescueSystem bool
		disk string
		serverNumber int
		version string
		image string
}

var initCmdFlags cmdFlags

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
			initializeServer(RobotClient, serverNumber, initCmdFlags)
	},
}

func init() {
		initCmd.Flags().BoolVarP(&initCmdFlags.skipReboot, "skipReboot", "n", false ,"skip reboot of server after enabling rescue system.")
		initCmd.Flags().BoolVarP(&initCmdFlags.enableRescueSystem, "enable-rescue-system", "r", false, "entering rescue system even if rescue system already enabled. This will generate a new password.")
		initCmd.Flags().StringVarP(&initCmdFlags.disk, "disk", "d", "nvme0n1", "disk to use for installation of image.")
		initCmd.Flags().StringVarP(&initCmdFlags.version, "version", "v", "1.9.2", "Talos version.")
		initCmd.Flags().StringVarP(&initCmdFlags.image, "image", "i", "", "Talos image URL. Don't use hcloud-amd64 image target Hetzner Cloud, use Talos 'metal' image instead.")
		addCommand(initCmd)
}

func initializeServer(client robot.Client, serverNumber int, f cmdFlags ) {
	sshPassword := os.Getenv("HETZNER_SSH_PASSWORD") // Set your Hetzner password in environment variable

	rescue, err := hetznerapi.GetRescueSystemDetails(client, serverNumber)
	if err != nil {
			fmt.Printf("Error getting rescue system status: %v\n", err)
			return
	}

	if (!rescue.Rescue.Active || f.enableRescueSystem) {
			rescue, err = hetznerapi.EnableRescueSystem(client, serverNumber)
			if err != nil {
					fmt.Printf("Error enabling rescue system: %v\n", err)
					return
			}
	}

	if !f.skipReboot  {
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

	version := "v1.9.2"
	if f.version != "" {
		if f.image != "" {
			fmt.Println("Warning: Both version and image flags are set. Using image flag.")
		}
		version = f.version
	}
	imageUrl := fmt.Sprintf("https://github.com/siderolabs/talos/releases/download/%s/metal-amd64.raw.zst", version)
	if f.image != "" {
		imageUrl = f.image
	}

	output, err := sshClient.DownloadImage(imageUrl)
	if err != nil {
		fmt.Printf("Failed to download image: %v, output %s\n", err, output)
		return
	}

	output, err = sshClient.InstallImage(f.disk)
	if err != nil {
		fmt.Printf("Failed to install image: %v output %s\n", err, output)
		return
	}

	hetznerapi.RebootServer(client, serverNumber)

	// Wait for Talos API to become available
	// Apply Talos configuration & bootstrap
	// Apply Cilium
	// Reboot node
}