package startrpc

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"../../pkg/common/config"
	"../../pkg/common/log"
	"../../pkg/common/mw"
	"../../pkg/common/network"
	"../../pkg/common/prome"
	"../../pkg/discoveryregistry"
	openKeeper "../../pkg/discoveryregistry/zookeeper"
	"../../pkg/utils"
	grpcPrometheus "../../libsrc/github.com/grpc-ecosystem/go-grpc-prometheus"
	"../../libsrc/google.golang.org/grpc"
	"../../libsrc/google.golang.org/grpc/credentials/insecure"
)

func Start(rpcPort int, rpcRegisterName string, /*prometheusPort int,*/ rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error, options ...grpc.ServerOption) error {
	fmt.Println("start", rpcRegisterName, "server, port: ", rpcPort, /*"prometheusPort:", prometheusPort,*/ ", ChatCore version: ", config.Version)
	listener, err := net.Listen("tcp", net.JoinHostPort(network.GetListenIP(config.Config.Rpc.ListenIP), strconv.Itoa(rpcPort)))
	if err != nil {
		return err
	}
	defer listener.Close()
	zkClient, err := openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
		openKeeper.WithFreq(time.Hour), openKeeper.WithUserNameAndPassword(config.Config.Zookeeper.Username,
			config.Config.Zookeeper.Password), openKeeper.WithRoundRobin(), openKeeper.WithTimeout(10), openKeeper.WithLogger(log.NewZkLogger()))
	if err != nil {
		return utils.Wrap1(err)
	}
	defer zkClient.CloseZK()
	zkClient.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	registerIP, err := network.GetRpcRegisterIP(config.Config.Rpc.RegisterIP)
	if err != nil {
		return err
	}
	/*
	// ctx 中间件
	if config.Config.Prometheus.Enable {
		prome.NewGrpcRequestCounter()
		prome.NewGrpcRequestFailedCounter()
		prome.NewGrpcRequestSuccessCounter()
		unaryInterceptor := mw.InterceptChain(grpcPrometheus.UnaryServerInterceptor, mw.RpcServerInterceptor)
		options = append(options, []grpc.ServerOption{
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(unaryInterceptor),
		}...)
	} else {
		options = append(options, mw.GrpcServer())
	}
	*/
	options = append(options, mw.GrpcServer())
	srv := grpc.NewServer(options...)
	defer srv.GracefulStop()
	err = rpcFn(zkClient, srv)
	if err != nil {
		return utils.Wrap1(err)
	}
	err = zkClient.Register(rpcRegisterName, registerIP, rpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	/*
	if err != nil {
		return utils.Wrap1(err)
	}
	go func() {
		if config.Config.Prometheus.Enable && prometheusPort != 0 {
			if err := prome.StartPrometheusSrv(prometheusPort); err != nil {
				panic(err.Error())
			}
		}
	}()
	*/
	return utils.Wrap1(srv.Serve(listener))
}
