package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
)

func main() {
	// Replace these variables with your actual SSH connection details
	sshHost := "localhost"
	sshPort := "22"
	sshUser := "your-ssh-username"
	sshPassword := "your-ssh-password"

	// SSH configuration
	sshConfig := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // You should handle host key verification in production
	}

	// SSH connection
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", sshHost, sshPort), sshConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer sshClient.Close()

	// Execute remote command
	command := "your-remote-command"
	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(command)
	if err != nil {
		log.Fatalf("Failed to run command '%s': %v\n%s", command, err, stderr.String())
	}

	fmt.Printf("Command output: %s", stdout.String())
}
