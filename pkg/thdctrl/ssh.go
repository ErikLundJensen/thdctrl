package thdctrl

import (
	"bytes"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

func EstablishSSHSession(host, port, user, password string) (*ssh.Session, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func ExecuteLSCommand(session *ssh.Session) (string, error) {
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("ls"); err != nil {
		return "", fmt.Errorf("failed to run command: %w", err)
	}
	return b.String(), nil
}

func DownloadImage(session *ssh.Session, url string) error {
	var b bytes.Buffer
	session.Stdout = &b
	download := fmt.Sprintf("wget -O /tmp/talos.raw.xz %s", url)
	if err := session.Run(download); err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	fmt.Printf("Downloaded image: %s", b.String())
	return nil
}

func InstallImage(session *ssh.Session, disk string) error {
	var b bytes.Buffer
	session.Stdout = &b
	unpack := fmt.Sprintf("zstdcat -dv /tmp/talos.raw.xz >/dev/%s", disk)
	if err := session.Run(unpack); err != nil {
		return fmt.Errorf("failed to run unpack image file: %w", err)
	}
	fmt.Printf("installed image at disk: %s", b.String())
	return nil
}

func WaitForReboot(sshHost string, sshPort string, sshUser string, sshPassword string) bool {
	maxRetries := 10
	retryInterval := 10 * time.Second

	for i := 0; i < maxRetries; i++ {
		fmt.Printf("Attempt %d: Establishing SSH session to %s:%s with user %s\n", i+1, sshHost, sshPort, sshUser)

		session, err := EstablishSSHSession(sshHost, sshPort, sshUser, sshPassword)
		if err != nil {
			fmt.Printf("Error establishing SSH session: %v\n", err)
			if i < maxRetries-1 {
				fmt.Printf("Retrying in %s...\n", retryInterval)
				time.Sleep(retryInterval)
				continue
			}
			return true
		}
		defer session.Close()
		return false
	}
	return true
}

