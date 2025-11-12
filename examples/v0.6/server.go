package main

import (
	"fmt"

	"github.com/lnnlmario/w2-inx/w2iface"
	"github.com/lnnlmario/w2-inx/w2net"
)

// ping test 自定义路由
type PingRouter struct {
	w2net.BaseRouter
}

// Ping Handle
func (this *PingRouter) Handle(request w2iface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().Send(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloW2inxRouter struct {
	w2net.BaseRouter
}

func (this *HelloW2inxRouter) Handle(request w2iface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().Send(1, []byte("Hello Zinx Router V0.6"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//创建一个server句柄
	s := w2net.NewServer()

	//配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloW2inxRouter{})

	//开启服务
	s.Serve()
}
