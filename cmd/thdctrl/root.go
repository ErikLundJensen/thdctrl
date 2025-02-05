package thdctrl

import (
	"os"
	"github.com/spf13/cobra"
	"github.com/eriklundjensen/thdctrl/pkg/robot"
)

var RobotClient = robot.Client {
	Username: os.Getenv("HETZNER_USERNAME"),
	Password: os.Getenv("HETZNER_PASSWORD"),
}

// Commands is a list of commands published by the package.
var Commands []*cobra.Command

func addCommand(cmd *cobra.Command) {
	Commands = append(Commands, cmd)
}

