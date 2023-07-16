package config

import (
	_ "embed"
)

//go:embed version
var Version string

var Config config

type CallBackConfig struct {
	Enable                 bool  `yaml:"enable"`
	CallbackTimeOut        int   `yaml:"timeout"`
	CallbackFailedContinue *bool `yaml:"failedContinue"`
}

type config struct {
	Zookeeper struct {
		Schema   string   `yaml:"schema"`
		ZkAddr   []string `yaml:"address"`
		Username string   `yaml:"username"`
		Password string   `yaml:"password"`
	} `yaml:"zookeeper"`

	Redis struct {
		Address  []string `yaml:"address"`
		Username string   `yaml:"username"`
		Password string   `yaml:"password"`
	} `yaml:"redis"`

	Kafka struct {
		Username         string   `yaml:"username"`
		Password         string   `yaml:"password"`
		Addr             []string `yaml:"addr"`
		LatestMsgToRedis struct {
			Topic string `yaml:"topic"`
		} `yaml:"latestMsgToRedis"`
		MsgToMongo struct {
			Topic string `yaml:"topic"`
		} `yaml:"offlineMsgToMongo"`
		MsgToPush struct {
			Topic string `yaml:"topic"`
		} `yaml:"msgToPush"`
		MsgToModify struct {
			Topic string `yaml:"topic"`
		} `yaml:"msgToModify"`
		ConsumerGroupID struct {
			MsgToRedis  string `yaml:"msgToRedis"`
			MsgToModify string `yaml:"msgToModify"`
		} `yaml:"consumerGroupID"`
	} `yaml:"kafka"`

	Rpc struct {
		RegisterIP string `yaml:"registerIP"`
		ListenIP   string `yaml:"listenIP"`
	} `yaml:"rpc"`

	Api struct {
		ChatCoreApiPort []int  `yaml:"chatCoreApiPort"`
		ListenIP      string `yaml:"listenIP"`
	} `yaml:"api"`

	RpcPort struct {
		ChatCoreUserPort           []int `yaml:"chatCoreUserPort"`
		ChatCoreMessagePort        []int `yaml:"chatCoreMessagePort"`
		ChatCoreMessageGatewayPort []int `yaml:"chatCoreMessageGatewayPort"`
	} `yaml:"rpcPort"`

	RpcRegisterName struct {
		ChatCoreUserName           string `yaml:"ChatCoreUserName"`
		ChatCoreMsgName            string `yaml:"ChatCoreMsgName"`
		ChatCoreMessageGatewayName string `yaml:"ChatCoreMessageGatewayName"`
	} `yaml:"rpcRegisterName"`

	Log struct {
		StorageLocation     string `yaml:"storageLocation"`
		RotationTime        int    `yaml:"rotationTime"`
		RemainRotationCount uint   `yaml:"remainRotationCount"`
		RemainLogLevel      int    `yaml:"remainLogLevel"`
		IsStdout            bool   `yaml:"isStdout"`
		IsJson              bool   `yaml:"isJson"`
		WithStack           bool   `yaml:"withStack"`
	} `yaml:"log"`

	LongConnSvr struct {
		OpenImWsPort        []int `yaml:"openImWsPort"`
		WebsocketMaxConnNum int   `yaml:"websocketMaxConnNum"`
		WebsocketMaxMsgLen  int   `yaml:"websocketMaxMsgLen"`
		WebsocketTimeout    int   `yaml:"websocketTimeout"`
	} `yaml:"longConnSvr"`

	MultiLoginPolicy                  int    `yaml:"multiLoginPolicy"`
	MsgCacheTimeout                   int    `yaml:"msgCacheTimeout"`
	Secret                            string `yaml:"secret"`
	TokenPolicy                       struct {
		Expire int64 `yaml:"expire"`
	} `yaml:"tokenPolicy"`

	Callback struct {
		CallbackUrl                        string         `yaml:"url"`
		CallbackBeforeSendSingleMsg        CallBackConfig `yaml:"beforeSendSingleMsg"`
		CallbackAfterSendSingleMsg         CallBackConfig `yaml:"afterSendSingleMsg"`
	} `yaml:"callback"`
}
