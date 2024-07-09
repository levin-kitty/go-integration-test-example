package tests

import (
	"encoding/json"
	sshServer "github.com/gliderlabs/ssh"
	"github.com/levin-kitty/go-test/app"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

// serverApiBaseUrl => don't need to assign port manually
// mockSSHServerAddr
// appBaseUrl

func TestApp_Uptime(t *testing.T) {
	// Arrange
	uptime := " 12:15:54 up 26 days,  5:45,  1 user,  load average: 0.00, 0.00, 0.00\n"

	sshServerAddr := "127.0.0.1:" + strconv.Itoa(randomPort())
	appAddr := "127.0.0.1:" + strconv.Itoa(randomPort())

	commandChan := make(chan string, 1)
	go RunMockSSHServer(t, sshServerAddr, func(session sshServer.Session) {
		commandChan <- session.Command()[0]

		if _, err := io.WriteString(session, uptime); err != nil {
			t.Error(err)
			return
		}
	})

	go func() {
		application, err := app.NewApp(".ssh/fake.key", "")
		if err != nil {
			t.Error(err)
			return
		}
		if err := application.Run(appAddr); err != nil {
			t.Error(err)
			return
		}
	}()

	waitServerReady(t, sshServerAddr, 1000, 10*time.Millisecond)
	waitServerReady(t, appAddr, 1000, 10*time.Millisecond)

	// Act
	resp, err := http.Get("http://" + appAddr + "/" + sshServerAddr + "/uptime")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	// Assert
	assert.Equal(t, uptime, string(body))
	select {
	case receivedCommand := <-commandChan:
		assert.Equal(t, "uptime", receivedCommand)
	case <-time.After(10 * time.Second):
		t.Error("command timeout")
	}
}

func TestApp_UptimeByServerId(t *testing.T) {
	// Arrange
	uptime := " 12:15:54 up 26 days,  5:45,  1 user,  load average: 0.00, 0.00, 0.00\n"

	sshServerAddr := "127.0.0.1:" + strconv.Itoa(randomPort())
	appAddr := "127.0.0.1:" + strconv.Itoa(randomPort())

	commandChan := make(chan string, 1)
	go RunMockSSHServer(t, sshServerAddr, func(session sshServer.Session) {
		commandChan <- session.Command()[0]

		if _, err := io.WriteString(session, uptime); err != nil {
			t.Error(err)
			return
		}
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/servers/{serverId}", func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(app.ServerResponse{IP: sshServerAddr})
		if err != nil {
			panic(err)
		}
		if _, err := w.Write(bytes); err != nil {
			panic(err)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	go func() {
		application, err := app.NewApp(".ssh/fake.key", server.URL)
		if err != nil {
			t.Error(err)
			return
		}
		if err = application.Run(appAddr); err != nil {
			t.Error(err)
			return
		}
	}()

	waitServerReady(t, sshServerAddr, 1000, 10*time.Millisecond)
	waitServerReady(t, appAddr, 1000, 10*time.Millisecond)

	// Act
	resp, err := http.Get("http://" + appAddr + "/servers/1/uptime")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	// Assert
	assert.Equal(t, uptime, string(body))
	select {
	case receivedCommand := <-commandChan:
		assert.Equal(t, "uptime", receivedCommand)
	case <-time.After(10 * time.Second):
		t.Error("command timeout")
	}
}
