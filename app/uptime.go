package app

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func (app *App) Uptime(w http.ResponseWriter, r *http.Request) {
	ip := r.PathValue("ip")
	output, err := app.SSHCommand(ip, "uptime")
	if err != nil {
		panic(err)
	}
	_, err = io.WriteString(w, output)
	if err != nil {
		panic(err)
	}
}

type ServerResponse struct {
	IP string `json:"ip"`
}

func (app *App) UptimeByServerId(w http.ResponseWriter, r *http.Request) {
	// get ip from server id
	serverId, err := strconv.Atoi(r.PathValue("serverId"))
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(app.serverApiConfig.baseUrl + "/servers/" + strconv.Itoa(serverId))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	server := new(ServerResponse)
	if err := json.Unmarshal(respBody, server); err != nil {
		panic(err)
	}

	// uptime
	output, err := app.SSHCommand(server.IP, "uptime")
	if err != nil {
		panic(err)
	}

	if _, err = io.WriteString(w, output); err != nil {
		panic(err)
	}
}
