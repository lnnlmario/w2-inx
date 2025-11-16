package w2iface

// 链接管理抽象层
type IConnManager interface {
	Add(conn IConnection)                   // 添加链接
	Remove(conn IConnection)                // 删除链接
	Get(connID uint32) (IConnection, error) // 获取链接
	Len() int                               // 当前链接个数
	ClearConn()                             // 删除并停止所有链接
}
