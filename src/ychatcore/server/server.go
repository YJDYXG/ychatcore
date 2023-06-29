package server

import (
	"net"
	"github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
	"github.com/Unknwon/goconfig"
	"os"
)

var (
	log = logrus.WithFields(logrus.Fields{"pkg": "service"})
)

type Handler interface {
	Init(c *goconfig.ConfigFile)
	Signal(sig os.Signal) bool //return true to stop daemon
}

type Service struct {
	*ServiceConfig
	handler      Handler
	tcpListener  net.Listener
	unixListener net.Listener
	Server       *grpc.Server
}

type ServiceConfig struct {
	Root  string
	Name  string
	Proto string
	Addr  string
	Sock  string
}

func (srv *Service) Start() {
	if srv.unixListener != nil {
		srv.Server.Serve(srv.unixListener)
	}

	if srv.tcpListener != nil {
		srv.Server.Serve(srv.tcpListener)
	}
}
