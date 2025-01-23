package main

import (
    "flag"
    "fmt"
    "github.com/eriklundjensen/thdctrl/pkg/thdctrl"
    "github.com/eriklundjensen/thdctrl/pkg/robot"
    "os"
)

func main() {
    command := flag.String("command", "listServers", "Command to run: [init] to initialize server, [listServers] to list servers, [listFirewallRules] to list firewall rules for a servers")
    serverNumber := flag.String("serverNumber", "2565351", "Server number")
    resetMode := flag.String("resetMode", "hw", "Reset mode: [hw] for hardware reset, [none] for no reset")
    enableRescueSystem := flag.String("enableRescueSystem", "true", "enableRescueSystem: [true] entering rescue system, [false] skip rescue system")
    disk := flag.String("disk", "nvme0n1", "Disk to install image")
    flag.Parse()

    var robotClient = robot.Client{
        Username: os.Getenv("HETZNER_USERNAME"),
        Password: os.Getenv("HETZNER_PASSWORD"),
    }

    switch(*command) {
    case "init":
        initializeServer(robotClient, resetMode, enableRescueSystem, disk, *serverNumber)
    case "listServers":
        listServers(robotClient)
    case "listFirewallRules":
        listFirewallRules(robotClient, *serverNumber)
    default:
        fmt.Printf("Command not recognized: %s\n", *command)
    }
}

func initializeServer(client robot.Client, resetMode *string, enableRescueSystem *string, disk *string, serverNumber string) {
    sshPassword := os.Getenv("HETZNER_SSH_PASSWORD") // Set your Hetzner password in environment variable

    rescue, err := thdctrl.GetRescueSystemDetails(client, serverNumber)
    if err != nil {
        fmt.Printf("Error getting rescue system status: %v\n", err)
        return
    }

    if (!rescue.Rescue.Active && *enableRescueSystem=="true") {
        rescue, err = thdctrl.EnableRescueSystem(client, serverNumber)
        if err != nil {
            fmt.Printf("Error enabling rescue system: %v\n", err)
            return
        }
    }

    if *resetMode != "none" {
        err = thdctrl.ResetServer(client, serverNumber, *resetMode)
    }
    if (err != nil || rescue == nil) {
        fmt.Printf("Rescue system state is not available: %v\n", err)
        return
    }
    sshHost := rescue.Rescue.ServerIP
    sshPort := "22" // Default SSH port
    sshUser := "root"
    if rescue.Rescue.Password != "" {
        sshPassword = rescue.Rescue.Password
    }

    thdctrl.WaitForReboot(sshHost, sshPort, sshUser, sshPassword)
	fmt.Printf("Server ready\n")

    // TODO: use wrapper for SSH session
	session, err := thdctrl.EstablishSSHSession(sshHost, sshPort, sshUser, sshPassword)
	if err != nil {
		fmt.Printf("Error establishing SSH session: %v\n", err)
		return
	}

    // Don't use image hcloud-amd64 target Hetzner Cloud, use Talos 'metal' image instead
    version := "v1.9.2"
    imageUrl := fmt.Sprintf("https://github.com/siderolabs/talos/releases/download/%s/metal-amd64.raw.zst", version)
    err = thdctrl.DownloadImage(session, imageUrl)
	if err != nil {
		fmt.Printf("Failed to download image: %v\n", err)
		return
	}
    session.Close()

    session, err = thdctrl.EstablishSSHSession(sshHost, sshPort, sshUser, sshPassword)
	if err != nil {
		fmt.Printf("Error establishing SSH session: %v\n", err)
		return
	}
    defer session.Close()

    err = thdctrl.InstallImage(session, *disk)
	if err != nil {
		fmt.Printf("Failed to install image: %v\n", err)
		return
	}

    // TODO: second reboot should not be dependen upon reset mode.
    if *resetMode != "none" {
        thdctrl.ResetServer(client, serverNumber, *resetMode)
    }

    // Wait for Talos API to become available
    // Apply Talos configuration
    // Apply Cilium
    // Reboot node
}

func listServers(client robot.Client) {
	servers, err := thdctrl.ListServers(client)
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

func listFirewallRules(client robot.Client, serverID string) {
    firewallRes, err := thdctrl.GetFirewallRules(client, serverID) 
    if err != nil {
        fmt.Printf("Error getting firewall rules: %v\n", err)
    }
    fmt.Println("Firewall status:")
    fmt.Printf("Server: %d, Status:%s\n", firewallRes.ServerNumber, firewallRes.Status)

    fmt.Println("Firewall rules:")
    for _, rule := range firewallRes.Rules {
        fmt.Printf("SrcIP: %s, DstIP: %s, Protocol: %s, SrcPort: %s, DstPort: %s, Action: %s, TCPFlags: %s\n",
            rule.SrcIP, rule.DstIP, rule.Protocol, rule.SrcPort, rule.DstPort, rule.Action, rule.TCPFlags)
    }
}