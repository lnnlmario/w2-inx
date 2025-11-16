package w2net

import (
	"errors"
	"fmt"
	"sync"

	"github.com/lnnlmario/w2-inx/w2iface"
)

// 链接管理模块
type ConnManager struct {
	connections map[uint32]w2iface.IConnection // 管理链接信息
	connLock    sync.RWMutex                   // 读写链接的读写锁
}

func (c *ConnManager) Add(conn w2iface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connections[conn.GetConnId()] = conn

	fmt.Println("connection add to ConnManager successfully: conn num=", c.Len())
}

func (c *ConnManager) Remove(conn w2iface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 删除链接信息
	delete(c.connections, conn.GetConnId())

	fmt.Println("connection remove from ConnManager successfully: conn num=", c.Len())
}

func (c *ConnManager) Get(connID uint32) (w2iface.IConnection, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	for connID, conn := range c.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(c.connections, connID)
	}

	fmt.Println("connection cleared successfully: conn num=", c.Len())
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]w2iface.IConnection),
	}
}
