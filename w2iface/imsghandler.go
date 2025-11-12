package w2iface

// 消息管理抽象层
type IMsgHandle interface {
	// 非阻塞方式处理消息
	DoMsgHandler(request IRequest)
	// 为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter)
}
