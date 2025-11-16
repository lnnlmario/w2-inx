package main

import (
	"fmt"

	"github.com/lnnlmario/w2-inx/w2iface"
	"github.com/lnnlmario/w2-inx/w2net"
)

type PingRouter struct {
	w2net.BaseRouter
}

// Ping Handle
func (this *PingRouter) Handle(request w2iface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloZinxRouter struct {
	w2net.BaseRouter
}

// HelloZinxRouter Handle
func (this *HelloZinxRouter) Handle(request w2iface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(1, []byte("Hello Zinx Router V0.8"))
	if err != nil {
		fmt.Println(err)
	}
}

// 创建连接的时候执行
func DoConnectionBegin(conn w2iface.IConnection) {
	fmt.Println("DoConnecionBegin is Called ... ")
	err := conn.Send(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

// 连接断开的时候执行
func DoConnectionLost(conn w2iface.IConnection) {
	fmt.Println("DoConneciotnLost is Called ... ")
}

func main() {
	s := w2net.NewServer()

	// 注册hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	// 开启服务
	s.Serve()
}
