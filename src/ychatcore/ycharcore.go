package ychatcore

import (
	"sync"
	"ychatcore/server"
	"github.com/Sirupsen/logrus"
)

var (
	log             = logrus.WithFields(logrus.Fields{"pkg": "ychatcore"})
	lastServicePids = make(map[uint32]int)
	rwLock          sync.RWMutex
)

type Ycharcore struct {
	srv            *service.Service
}

func New(conf, name, proto, addr string) (*Ycharcore, error) {

}

func (a *Agent) Start() {
	ichatcore.RegisterChatCoreServer(a.Service.Server, a)
	ichatcore.RegisterNetworkServiceServer(a.Service.Server, a)
	a.Service.Start()
}
