package msggateway

import (
	"fmt"
	"time"

	"../../pkg/common/config"
)

func RunWsAndServer(rpcPort, wsPort, prometheusPort int) error {
	fmt.Println("start rpc/msg_gateway server, port: ", rpcPort, wsPort, prometheusPort, ", OpenIM version: ", config.Version)
	longServer, err := NewWsServer(
		WithPort(wsPort),
		WithMaxConnNum(int64(config.Config.LongConnSvr.WebsocketMaxConnNum)),
		WithHandshakeTimeout(time.Duration(config.Config.LongConnSvr.WebsocketTimeout)*time.Second),
		WithMessageMaxMsgLength(config.Config.LongConnSvr.WebsocketMaxMsgLen))
	if err != nil {
		return err
	}
	hubServer := NewServer(rpcPort, longServer)
	go hubServer.Start()
	hubServer.LongConnServer.Run()
	return nil
}

