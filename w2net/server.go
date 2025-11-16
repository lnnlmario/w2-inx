package w2net

import (
	"fmt"
	"net"

	"github.com/lnnlmario/w2-inx/w2utils"

	"github.com/lnnlmario/w2-inx/w2iface"
)

type Server struct {
	// 服务器的而名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器监听的IP
	IP string
	// 服务器监听的端口
	Port int
	// 消息管理模块，绑定msgId和对应的处理的方法
	msgHandler w2iface.IMsgHandle
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server is listening IP:%s, Port:%d, Version:%s\n", s.IP, s.Port, s.IPVersion)

	fmt.Printf("[W2inx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		w2utils.GlobalObject.Version,
		w2utils.GlobalObject.MaxConn,
		w2utils.GlobalObject.MaxPacketSize)

	go func() {
		// 0. 启动worker工作池机制
		s.msgHandler.StartWorkerPool()

		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("ResolveTCPAddr error:", err)
			return
		}

		// 2. 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("ListenTCP error:", err)
			return
		}
		fmt.Println("start w2-inx server", s.Name, "success, now listening...")

		var cId uint32 = 0

		// 3. 阻塞等待客户端的链接，处理客户端链接业务(读写)
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept error:", err)
				continue
			}

			// 将处理新链接的业务方法和conn进行绑定，得到我们的链接模块
			dealConn := NewConnection(conn, cId, s.msgHandler)
			cId++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] w2-inx server , name ", s.Name)
}

func (s *Server) Serve() {
	// 启动服务器
	s.Start()

	// 阻塞主线程
	select {}
}

func (s *Server) AddRouter(msgId uint32, router w2iface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

/**
 * NewServer 创建一个服务器实例
 * @param name 服务器名称
 * @return IServer 服务器实例
 */
func NewServer(args ...string) w2iface.IServer {
	serverName := w2utils.GlobalObject.Name
	if len(args) != 0 {
		serverName = args[0]
	}

	s := &Server{
		Name:       serverName,
		IPVersion:  "tcp4",
		IP:         w2utils.GlobalObject.Host,
		Port:       w2utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
	}

	return s
}
