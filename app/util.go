package app

import (
	"bytes"
	"golang.org/x/crypto/ssh"
)

func (app *App) SSHCommand(addr string, cmd string) (string, error) {
	client, err := ssh.Dial("tcp", addr, app.sshConfig.sshClientConfig)
	if err != nil {
		return "", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	if err := session.Run(cmd); err != nil {
		return "", err
	}
	return buf.String(), nil
}
