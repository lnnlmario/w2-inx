package w2net

import (
	"errors"
	"fmt"
	"github.com/lnnlmario/w2-inx/w2utils"
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

	// 无缓冲管道 用于读写两个goroutine间消息
	msgChan chan []byte
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connId uint32, msgHandler w2iface.IMsgHandle) *Connection {
	return &Connection{
		Conn:       conn,
		ConnId:     connId,
		isClosed:   false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
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

		// 判断用户配置的workerPoolSize个数，若<=0，走之前启动一个客户端处理请求消息
		if w2utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经启动工作池机制，将消息交给worker处理
			c.MsgHandler.SendMsgToTaskQueue(&request)
		} else {
			// 从绑定好的消息和对应的处理方法中执行对应的handle方法
			go c.MsgHandler.DoMsgHandler(&request)
		}
	}
}

func (c *Connection) StartWriter() {
	fmt.Println("Write goroutine is running..., connId=", c.ConnId)
	defer fmt.Println("Write connId=", c.ConnId, "remote addr=", c.RemoteAddr(), ", conn writer exit!")

	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data err:", err, ", connId=", c.ConnId, ", conn writer exit!")
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

// 启动链接 让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Connection Start()..., connId=", c.ConnId)
	// 启动从当前链接读数据的业务
	go c.StartReader()
	// 启动从当前链接写数据的业务
	go c.StartWriter()

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
	close(c.msgChan)
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
	c.msgChan <- msg

	return nil
}
