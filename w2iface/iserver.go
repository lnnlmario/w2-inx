package w2iface

// 定义一个服务器接口
type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()
	// 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	AddRouter(msgId uint32, router IRouter)
	// 得到链接管理器
	GetConnMgr() IConnManager
	// 设置连接创建前的Hook函数
	SetOnConnStart(func(connection IConnection))
	// 设置连接断开时的Hook函数
	SetOnConnStop(func(connection IConnection))
	// 调用连接开始前 Hook函数
	CallOnConnStart(connection IConnection)
	// 调用连接断开时 Hook函数
	CallOnConnStop(connection IConnection)
}
