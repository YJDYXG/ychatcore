package msggateway

import (
	"context"
	"../../pkg/common/db/cache"

	"../../pkg/common/config"
	"../../pkg/common/constant"
	"../../pkg/common/log"
	"../../common/prome"
	"../../common/tokenverify"
	"../../discoveryregistry"
	"../../pkg/errs"
	"../../pkg/proto/msggateway"
	"../../pkg/startrpc"
	"../../pkg/utils"
	"../../libsrc/google.golang.org/grpc"
)

func (s *Server) InitServer(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	msgModel := cache.NewMsgCacheModel(rdb)
	s.LongConnServer.SetDiscoveryRegistry(client)
	s.LongConnServer.SetCacheHandler(msgModel)
	msggateway.RegisterMsgGatewayServer(server, s)
	return nil
}

func (s *Server) Start() error {
	return startrpc.Start(s.rpcPort, config.Config.RpcRegisterName.ChatCoreMessageGatewayName,/* s.prometheusPort,*/ s.InitServer)
}

type Server struct {
	rpcPort        int
	prometheusPort int
	LongConnServer LongConnServer
	pushTerminal   []int
}

func (s *Server) SetLongConnServer(LongConnServer LongConnServer) {
	s.LongConnServer = LongConnServer
}

func NewServer(rpcPort int, longConnServer LongConnServer) *Server {
	return &Server{rpcPort: rpcPort, LongConnServer: longConnServer, pushTerminal: []int{constant.IOSPlatformID, constant.AndroidPlatformID}}
}

func (s *Server) OnlinePushMsg(context context.Context, req *msggateway.OnlinePushMsgReq) (*msggateway.OnlinePushMsgResp, error) {
	panic("implement me")
}

func (s *Server) GetUsersOnlineStatus(ctx context.Context, req *msggateway.GetUsersOnlineStatusReq) (*msggateway.GetUsersOnlineStatusResp, error) {
	if !tokenverify.IsAppManagerUid(ctx) {
		return nil, errs.ErrNoPermission.Wrap("only app manager")
	}
	var resp msggateway.GetUsersOnlineStatusResp
	for _, userID := range req.UserIDs {
		clients, ok := s.LongConnServer.GetUserAllCons(userID)
		if !ok {
			continue
		}
		temp := new(msggateway.GetUsersOnlineStatusResp_SuccessResult)
		temp.UserID = userID
		for _, client := range clients {
			if client != nil {
				ps := new(msggateway.GetUsersOnlineStatusResp_SuccessDetail)
				ps.Platform = constant.PlatformIDToName(client.PlatformID)
				ps.Status = constant.OnlineStatus
				ps.ConnID = client.ctx.GetConnID()
				ps.IsBackground = client.IsBackground
				temp.Status = constant.OnlineStatus
				temp.DetailPlatformStatus = append(temp.DetailPlatformStatus, ps)
			}
		}
		if temp.Status == constant.OnlineStatus {
			resp.SuccessResult = append(resp.SuccessResult, temp)
		}
	}
	return &resp, nil
}

func (s *Server) OnlineBatchPushOneMsg(ctx context.Context, req *msggateway.OnlineBatchPushOneMsgReq) (*msggateway.OnlineBatchPushOneMsgResp, error) {
	panic("implement me")
}


func (s *Server) KickUserOffline(ctx context.Context, req *msggateway.KickUserOfflineReq) (*msggateway.KickUserOfflineResp, error) {
	for _, v := range req.KickUserIDList {
		if clients, _, ok := s.LongConnServer.GetUserPlatformCons(v, int(req.PlatformID)); ok {
			for _, client := range clients {
				err := client.KickOnlineMessage()
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return &msggateway.KickUserOfflineResp{}, nil
}

func (s *Server) MultiTerminalLoginCheck(ctx context.Context, req *msggateway.MultiTerminalLoginCheckReq) (*msggateway.MultiTerminalLoginCheckResp, error) {
	//TODO implement me
	panic("implement me")
}
