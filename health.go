package main

import (
	"context"
	"fmt"
	"net"
	"strings"
)

const tcpSocket = "127.0.0.1:9999"
const BUF_SIZE = 4096

type HealthcheckServer struct {
	socket          net.Listener
	ctx             context.Context
	cancel          context.CancelFunc
	healthcheckFunc func() error
}

// Init initializes the HealthcheckServer instance by setting the socketPath
// and creating a new context with cancel function.
//
// Example:
//
//	srv := HealthcheckServer{}
//	srv.Init()
func (s *HealthcheckServer) Init(f func() error) {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.healthcheckFunc = f
}

func (s *HealthcheckServer) Start() (err error) {
	s.socket, err = net.Listen("tcp", tcpSocket)
	if err != nil {
		return err
	}
	for {
		select {
		case <-s.ctx.Done():
			_ = s.socket.Close()
			return nil
		default:
			conn, err := s.socket.Accept()
			if err != nil {
				continue
			}
			go func(conn net.Conn) {
				defer conn.Close()
				inputBuffer := make([]byte, BUF_SIZE)
				n, _ := conn.Read(inputBuffer)
				incomingMessage := strings.TrimSpace(string(inputBuffer[:n-1]))
				if incomingMessage != "ping" {
					_, _ = conn.Write([]byte("send 'ping' to trigger healthcheck\n"))
					return
				}
				err = s.healthcheckFunc()
				if err != nil {
					conn.Write([]byte(fmt.Sprintf("%s", err.Error())))
					return
				}
				conn.Write([]byte("success\n"))
			}(conn)
		}
	}
}

func (s *HealthcheckServer) Stop() {
	s.cancel()
}
