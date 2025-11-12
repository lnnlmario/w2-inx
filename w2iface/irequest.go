package w2iface

// IRequest 定义一个请求接口
// 实际上是把客户端请求的链接信息和请求数据包装到一起，包装客户端的全部请求数据
type IRequest interface {
	// /获取请求连接信息
	GetConnection() IConnection
	// 获取请求消息的数据
	GetData() []byte
	// 获取请求消息的ID
	GetMsgID() uint32
}
