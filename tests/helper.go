package tests

import (
	"crypto/rand"
	sshServer "github.com/gliderlabs/ssh"
	"math/big"
	"net"
	"testing"
	"time"
)

func randomPort() int {
	n, err := rand.Int(rand.Reader, big.NewInt(30001))
	if err != nil {
		panic(err)
	}
	return int(n.Int64()) + 10000
}

func waitServerReady(t *testing.T, address string, maxRetries int, delay time.Duration) {
	for i := 0; i < maxRetries; i++ {
		conn, err := net.Dial("tcp", address)
		if err == nil {
			err := conn.Close()
			if err != nil {
				t.Error(err)
			}
			return
		}
		time.Sleep(delay)
	}
	t.Errorf("server startup got timeout: address=%s", address)
}

/* Mock SSH Server */

type MockSSHServer struct {
	sshServer sshServer.Server
}

func (m *MockSSHServer) Handle(handler func(session sshServer.Session)) {
	m.sshServer.Handle(handler)
}

func (m *MockSSHServer) ListenAndServe() error {
	return m.sshServer.ListenAndServe()
}

func NewMockSSHServer(addr string) *MockSSHServer {
	return &MockSSHServer{sshServer: sshServer.Server{
		Addr: addr,
	}}
}

func RunMockSSHServer(t *testing.T, addr string, handler func(sshServer.Session)) {
	mock := NewMockSSHServer(addr)
	mock.Handle(handler)
	if err := mock.ListenAndServe(); err != nil {
		t.Fatal(err)
	}
}
