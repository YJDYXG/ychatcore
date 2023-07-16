package api

import (
	"context"

	"../../pkg/common/config"
	"../..//pkg/common/log"
	"../..//pkg/common/mw"
	"../..//pkg/common/prome"
	"../..//pkg/discoveryregistry"
	"../../libsrc/github.com/gin-gonic/gin"
	"../../libsrc/github.com/gin-gonic/gin/binding"
	"../../libsrc/github.com/go-playground/validator/v10"
	"../../libsrc/github.com/redis/go-redis/v9"
	"../../libsrc/google.golang.org/grpc"
	"../../libsrc/google.golang.org/grpc/credentials/insecure"
)

func NewGinRouter(discov discoveryregistry.SvcDiscoveryRegistry, rdb redis.UniversalClient) *gin.Engine {
	discov.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials())) // 默认RPC中间件
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("required_if", RequiredIf)
	}
	log.ZInfo(context.Background(), "load config", "config", config.Config)
	r.Use(gin.Recovery(), mw.CorsHandler(), mw.GinParseOperationID())
	u := NewUserApi(discov)
	/*
	if config.Config.Prometheus.Enable {
		prome.NewApiRequestCounter()
		prome.NewApiRequestFailedCounter()
		prome.NewApiRequestSuccessCounter()
		r.Use(prome.PrometheusMiddleware)
		r.GET("/metrics", prome.PrometheusHandler())
	}
	*/

	ParseToken := mw.GinParseToken(rdb)
	userRouterGroup := r.Group("/user")
	{
		userRouterGroup.POST("/user_register", u.UserRegister)
		userRouterGroup.POST("/update_user_info", ParseToken, u.UpdateUserInfo)
		userRouterGroup.POST("/set_global_msg_recv_opt", ParseToken, u.SetGlobalRecvMessageOpt)
		userRouterGroup.POST("/get_users_info", ParseToken, u.GetUsersPublicInfo)
		userRouterGroup.POST("/get_all_users_uid", ParseToken, u.GetAllUsersID)
		userRouterGroup.POST("/account_check", ParseToken, u.AccountCheck)
		userRouterGroup.POST("/get_users", ParseToken, u.GetUsers)
		userRouterGroup.POST("/get_users_online_status", ParseToken, u.GetUsersOnlineStatus)
	}
	//certificate
	authRouterGroup := r.Group("/auth")
	{
		a := NewAuthApi(discov)
		authRouterGroup.POST("/user_register", u.UserRegister)
		authRouterGroup.POST("/user_token", a.UserToken)
		authRouterGroup.POST("/parse_token", a.ParseToken)
		authRouterGroup.POST("/force_logout", ParseToken, a.ForceLogout)
	}
	//Message
	msgGroup := r.Group("/msg", ParseToken)
	{
		m := NewMessageApi(discov)
		msgGroup.POST("/newest_seq", m.GetSeq)
		msgGroup.POST("/send_msg", m.SendMessage)
		msgGroup.POST("/pull_msg_by_seq", m.PullMsgBySeqs)
		msgGroup.POST("/revoke_msg", m.RevokeMsg)
		msgGroup.POST("/mark_msgs_as_read", m.MarkMsgsAsRead)
		msgGroup.POST("/mark_conversation_as_read", m.MarkConversationAsRead)
		msgGroup.POST("/get_conversations_has_read_and_max_seq", m.GetConversationsHasReadAndMaxSeq)
		msgGroup.POST("/set_conversation_has_read_seq", m.SetConversationHasReadSeq)

		msgGroup.POST("/clear_conversation_msg", m.ClearConversationsMsg)
		msgGroup.POST("/user_clear_all_msg", m.UserClearAllMsg)
		msgGroup.POST("/delete_msgs", m.DeleteMsgs)
		msgGroup.POST("/delete_msg_phsical_by_seq", m.DeleteMsgPhysicalBySeq)
		msgGroup.POST("/delete_msg_physical", m.DeleteMsgPhysical)

		msgGroup.POST("/batch_send_msg", m.ManagementBatchSendMsg)
		msgGroup.POST("/check_msg_is_send_success", m.CheckMsgIsSendSuccess)
	}
	//Conversation
	conversationGroup := r.Group("/conversation", ParseToken)
	{
		c := NewConversationApi(discov)
		conversationGroup.POST("/get_all_conversations", c.GetAllConversations)
		conversationGroup.POST("/get_conversation", c.GetConversation)
		conversationGroup.POST("/get_conversations", c.GetConversations)
		conversationGroup.POST("/batch_set_conversation", c.BatchSetConversations)
		conversationGroup.POST("/set_recv_msg_opt", c.SetRecvMsgOpt)
		conversationGroup.POST("/modify_conversation_field", c.ModifyConversationField)
		conversationGroup.POST("/set_conversations", c.SetConversations)
	}

	statisticsGroup := r.Group("/statistics", ParseToken)
	{
		statisticsGroup.POST("/user_register", u.UserRegisterCount)
	}
	return r
}
