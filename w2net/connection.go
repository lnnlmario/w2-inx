package w2net

import (
	"fmt"
	"net"

	"github.com/lnnlmario/w2-inx/w2iface"
)

type Connection struct {
	//当前链接的socker TCP套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前的链接状态
	isClosed bool

	//当前链接所绑定的处理业务方法API
	HandleAPI w2iface.HandleFunc

	//告知链接已经退出/停止的 channel
	ExitChan chan bool
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connId = ", c.ConnID, " Reader is exit, remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中，最大512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("connId = ", c.ConnID, " recv buf err:", err)
			continue
		}

		// 调用链接所绑定的HandleAPI
		if err := c.HandleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("connId = ", c.ConnID, " handle is error:", err)
			break
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conntion stop()... ConnId = ", c.ConnID)

	// 启动从当前链接读数据的业务
	go c.StartReader()

	// TODO: 启动从当前链接写数据的业务
}

func (c *Connection) Stop() {
	fmt.Println("Conntion stop()... ConnId = ", c.ConnID)

	if c.isClosed {
		return
	}
	c.isClosed = true

	// 关闭socket链接
	c.Conn.Close()

	// 回收资源
	close(c.ExitChan)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	//TODO implement me
	panic("implement me")
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, callback w2iface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		HandleAPI: callback,
		ExitChan:  make(chan bool, 1),
	}

	return c
}
