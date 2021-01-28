package server

import (
	"demo/manager"
	"net"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

type Server struct {
	logger   log.Logger
	listener net.Listener
	manager  *manager.Manager
}

func New(logger log.Logger, network, addr string, manager *manager.Manager) (*Server, error) {
	var (
		err    error
		server = &Server{
			logger:  logger,
			manager: manager,
		}
	)
	server.listener, err = net.Listen(network, addr)
	if err != nil {
		return nil, errors.Wrap(err, "server: create server error")
	}
	level.Info(server.logger).Log("msg", "TCP server start", "address", addr)
	return server, nil
}

func (s *Server) Close() {
	level.Info(s.logger).Log("msg", "Server close")
	s.listener.Close()
}

func (s *Server) Run() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return errors.Wrap(err, "server: accept connection error")
		}
		level.Info(s.logger).Log("msg", "Client connect", "addr", conn.RemoteAddr())
		go s.manager.Connect(conn)
	}
}
