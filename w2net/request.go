package w2net

import "github.com/lnnlmario/w2-inx/w2iface"

type Request struct {
	// 和客户端建立好的链接
	conn w2iface.IConnection
	// 客户端请求的数据
	data []byte
}

// 得到当前链接
func (r *Request) GetConnection() w2iface.IConnection {
	return r.conn
}

// 得到请求的消息数据
func (r *Request) GetData() []byte {
	return r.data
}
