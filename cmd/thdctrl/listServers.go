package thdctrl

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/eriklundjensen/thdctrl/pkg/robot"
	"github.com/eriklundjensen/thdctrl/pkg/hetznerapi"
)

var listServersCmd = &cobra.Command{
	Use:   "listServers",
	Short: "List all servers",
	Run: func(cmd *cobra.Command, args []string) {
		listServers(RobotClient)
	},
}

func init() {
	addCommand(listServersCmd)
}

func listServers(client robot.Client) {
	servers, err := hetznerapi.ListServers(client)
	if err != nil {
		fmt.Printf("Error listing servers: %v\n", err)
	}
	fmt.Println("List of servers:")
	for _, server := range servers {
		serverDetail := server.Server
		fmt.Printf("ID: %d, Name: %s, Product: %s, Datacenter: %s, IPv4: %s, IPv6: %s\n",
			serverDetail.ServerNumber, serverDetail.ServerName, serverDetail.Product, serverDetail.Datacenter, serverDetail.ServerIP, serverDetail.ServerIPv6Net)
	}
}
