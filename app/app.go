package app

import (
	"golang.org/x/crypto/ssh"
	"net/http"
	"os"
)

type SSHConfig struct {
	sshClientConfig *ssh.ClientConfig
}

type ServerApiConfig struct {
	baseUrl string
}

type App struct {
	sshConfig       *SSHConfig
	serverApiConfig *ServerApiConfig
}

func NewApp(privateKeyPath string, serverApiBaseUrl string) (*App, error) {
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	key, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	sshClientConfig := ssh.ClientConfig{
		User: "opc",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConfig := &SSHConfig{
		sshClientConfig: &sshClientConfig,
	}
	return &App{
		sshConfig: sshConfig,
		serverApiConfig: &ServerApiConfig{
			baseUrl: serverApiBaseUrl,
		},
	}, nil
}

func (app *App) Run(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/{ip}/uptime", app.Uptime)
	mux.HandleFunc("/servers/{serverId}/uptime", app.UptimeByServerId)
	return http.ListenAndServe(addr, mux)
}
