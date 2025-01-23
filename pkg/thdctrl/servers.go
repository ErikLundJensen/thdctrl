package thdctrl

import (
	"encoding/json"
	"fmt"
	"net/url"
	"github.com/eriklundjensen/thdctrl/pkg/robot"
)

type Subnet struct {
	IP   string `json:"ip"`
	Mask string `json:"mask"`
}

type ServerDetails struct {
	ServerIP         string      `json:"server_ip"`
	ServerIPv6Net    string      `json:"server_ipv6_net"`
	ServerNumber     int         `json:"server_number"`
	ServerName       string      `json:"server_name"`
	Product          string      `json:"product"`
	Datacenter       string      `json:"dc"`
	Traffic          string      `json:"traffic"`
	Status           string      `json:"status"`
	Cancelled        bool        `json:"cancelled"`
	PaidUntil        string      `json:"paid_until"`
	IP               []string    `json:"ip"`
	Subnet           []Subnet    `json:"subnet"`
	LinkedStorageBox interface{} `json:"linked_storagebox"`
}

type Server struct {
	Server ServerDetails `json:"server"`
}

type Servers struct {
	Servers Server `json:"server"`
}

func ListServers(client robot.Client) ([]Server, error) {
	path := "server"
	
	body, err := client.Get(path)
	if err != nil {
		return nil, err
	}

	var servers []Server
	if err := json.Unmarshal(body, &servers); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return servers, nil
}

func ResetServer(client robot.Client, serverNumber string, resetMode string) error {
	path := fmt.Sprintf("reset/%s", serverNumber)

	data := url.Values{}
	data.Set("type", resetMode)

	_, err := client.Post(path, data)
	if err != nil {
		return err
	}
	fmt.Println("Server reset successfully.")
	return nil
}
