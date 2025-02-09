package hetznerapi

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

func ListServers(client robot.ClientInterface) ([]Server, error) {
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

func RebootServer(client robot.ClientInterface, serverNumber int) error {
	path := fmt.Sprintf("reset/%d", serverNumber)

	data := url.Values{}
	data.Set("type", "hw")

	_, err := client.Post(path, data)
	if err != nil {
		return err
	}
	fmt.Println("Server reboot successfully initiated")
	return nil
}
