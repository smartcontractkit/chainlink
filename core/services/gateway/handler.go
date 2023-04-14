package gateway

type Callback interface {
	SendResponse(msg *Message)
}

type Handler interface {
	Init(connMgr ConnectionManager, donConfig *GatewayDONConfig)

	HandleUserMessage(msg *Message, cb Callback)

	HandleNodeMessage(msg *Message, nodeAddr string)
}
