package msggateway

import (
	"context"
	cbapi "../../pkg/callbackstruct"
	"../../pkg/common/config"
	"../../pkg/common/constant"
	"../../pkg/common/http"
	"../../pkg/common/mcontext"
	"time"
)

func url() string {
	return config.Config.Callback.CallbackUrl
}

func CallbackUserOnline(ctx context.Context, userID string, platformID int, isAppBackground bool, connID string) error {
	if !config.Config.Callback.CallbackUserOnline.Enable {
		return nil
	}
	req := cbapi.CallbackUserOnlineReq{
		UserStatusCallbackReq: cbapi.UserStatusCallbackReq{
			UserStatusBaseCallback: cbapi.UserStatusBaseCallback{
				CallbackCommand: constant.CallbackUserOnlineCommand,
				OperationID:     mcontext.GetOperationID(ctx),
				PlatformID:      platformID,
				Platform:        constant.PlatformIDToName(platformID),
			},
			UserID: userID,
		},
		Seq:             time.Now().UnixMilli(),
		IsAppBackground: isAppBackground,
		ConnID:          connID,
	}
	resp := cbapi.CommonCallbackResp{}
	return http.CallBackPostReturn(ctx, url(), &req, &resp, config.Config.Callback.CallbackUserOnline)
}

func CallbackUserOffline(ctx context.Context, userID string, platformID int, connID string) error {
	if !config.Config.Callback.CallbackUserOffline.Enable {
		return nil
	}
	req := &cbapi.CallbackUserOfflineReq{
		UserStatusCallbackReq: cbapi.UserStatusCallbackReq{
			UserStatusBaseCallback: cbapi.UserStatusBaseCallback{
				CallbackCommand: constant.CallbackUserOfflineCommand,
				OperationID:     mcontext.GetOperationID(ctx),
				PlatformID:      platformID,
				Platform:        constant.PlatformIDToName(platformID),
			},
			UserID: userID,
		},
		Seq:    time.Now().UnixMilli(),
		ConnID: connID,
	}
	resp := &cbapi.CallbackUserOfflineResp{}
	return http.CallBackPostReturn(ctx, url(), req, resp, config.Config.Callback.CallbackUserOffline)
}

func CallbackUserKickOff(ctx context.Context, userID string, platformID int) error {
	if !config.Config.Callback.CallbackUserKickOff.Enable {
		return nil
	}
	req := &cbapi.CallbackUserKickOffReq{
		UserStatusCallbackReq: cbapi.UserStatusCallbackReq{
			UserStatusBaseCallback: cbapi.UserStatusBaseCallback{
				CallbackCommand: constant.CallbackUserKickOffCommand,
				OperationID:     mcontext.GetOperationID(ctx),
				PlatformID:      platformID,
				Platform:        constant.PlatformIDToName(platformID),
			},
			UserID: userID,
		},
		Seq: time.Now().UnixMilli(),
	}
	resp := &cbapi.CommonCallbackResp{}
	return http.CallBackPostReturn(ctx, url(), req, resp, config.Config.Callback.CallbackUserOffline)
}

