package w2net

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/lnnlmario/w2-inx/w2iface"
)

// 链接模块
type Connection struct {
	// 当前链接的socket TCP套接字
	Conn *net.TCPConn

	// 链接的ID
	ConnId uint32

	// 当前的链接状态
	isClosed bool

	// 消息管理MsgId和对应处理方法的消息管理模块
	MsgHandler w2iface.IMsgHandle

	// 告知当前链接已经退出/停止 channel
	ExitChan chan bool
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connId uint32, msgHandler w2iface.IMsgHandle) *Connection {
	return &Connection{
		Conn:       conn,
		ConnId:     connId,
		isClosed:   false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
	}
}

func (c *Connection) StartReader() {
	fmt.Println("Read goroutine is running..., connId=", c.ConnId)
	defer fmt.Println("Read connId=", c.ConnId, "remote addr=", c.RemoteAddr(), "is exit!")
	defer c.Stop()
	// 循环读取数据
	for {
		// 创建拆包解包对象
		dp := NewDataPack()

		// 读取客户端你的msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head data err:", err)
			c.ExitChan <- true
			continue
		}

		// 拆包得到 msgid 和 datalen
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack msg err:", err)
			c.ExitChan <- true
			continue
		}

		// 根据 datalen 读取 data
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err:", err)
				c.ExitChan <- true
				continue
			}
		}
		msg.SetData(data)

		// 得到当前客户端请求的Request数据
		request := Request{conn: c, msg: msg}
		// 从绑定好的消息和对应的处理方法中执行对应的handle方法
		go c.MsgHandler.DoMsgHandler(&request)
	}
}

// 启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Connection Start()..., connId=", c.ConnId)
	// 启动从当前链接读数据的业务
	go c.StartReader()
	// TODO: 启动从当前链接写数据的业务

	for {
		select {
		case <-c.ExitChan:
			//得到退出消息，不再阻塞
			return
		}
	}
}

// 停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Connection Stop()..., connId=", c.ConnId)

	// 如果当前链接已经关闭
	if c.isClosed {
		return
	}
	// 将当前链接状态设置为已关闭
	c.isClosed = true
	// 关闭socket链接
	if err := c.Conn.Close(); err != nil {
		fmt.Println("Connection Stop() Close() error:", err)
		return
	}
	// 关闭socket链接
	c.Conn.Close()
	// 通知当前链接已经退出
	c.ExitChan <- true
	// 关闭当前链接的ExitChan
	close(c.ExitChan)
}

// 获取当前链接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前链接模块的链接ID
func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

// 获取远程客户端的TCP状态地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据 将数据发送给远程的客户端
func (c *Connection) Send(msgId uint32, data []byte) error {
	if c.isClosed {
		return errors.New("connection is closed")
	}
	// 将data封包发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack msg error:", err, "msgId=", msgId)
		return err
	}
	// 写回客户端
	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("send msg error:", err, "msgId=", msgId)
		c.ExitChan <- true
		return err
	}
	return nil
}
