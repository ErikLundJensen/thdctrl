package thdctrl

import (
	"bytes"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	Host, Port string
	Session *ssh.Session
	Config *ssh.ClientConfig
}

func (client *SSHClient) Auth(user, password string) {
	client.Config = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
}

func (client *SSHClient) establishSSHSession() error {
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", client.Host, client.Port), client.Config)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	client.Session = session
	return nil
}

func (client *SSHClient) ExecuteCommand(command string) (string, error) {
	var b bytes.Buffer
	client.Session.Stdout = &b
	client.establishSSHSession()
	defer client.Session.Close()
	
	if err := client.Session.Run(command); err != nil {
		return "", fmt.Errorf("failed to run command: %w", err)
	}
	return b.String(), nil
}

func (client *SSHClient) ExecuteLSCommand() (string, error) {
	return client.ExecuteCommand("ls")
}

func (client *SSHClient) DownloadImage(url string) (string,error) {
	download := fmt.Sprintf("wget -O /tmp/talos.raw.xz %s", url)
	return client.ExecuteCommand(download)
}

func (client *SSHClient) InstallImage(disk string) (string,error) {
	unpack := fmt.Sprintf("zstdcat -dv /tmp/talos.raw.xz >/dev/%s", disk)
	return client.ExecuteCommand(unpack)
}

func (client *SSHClient) WaitForReboot() bool {
	maxRetries := 10
	retryInterval := 10 * time.Second

	for i := 0; i < maxRetries; i++ {
		fmt.Printf("Attempt %d: Establishing SSH session to %s:%s\n", i+1, client.Host, client.Port)

		err := client.establishSSHSession()
		if err != nil {
			fmt.Printf("Error establishing SSH session: %v\n", err)
			if i < maxRetries-1 {
				fmt.Printf("Retrying in %s...\n", retryInterval)
				time.Sleep(retryInterval)
				continue
			}
			return true
		}
		return false
	}
	return true
}

