package main

import (
	"fmt"

	"github.com/lnnlmario/w2-inx/w2iface"
	"github.com/lnnlmario/w2-inx/w2net"
)

type PingRouter struct {
	w2net.BaseRouter
}

func (this *PingRouter) Handle(request w2iface.IRequest) {
	fmt.Println("PingRouter Handle")

	// 读取客户端的数据
	fmt.Println("receive from client: msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	// 回写数据
	err := request.GetConnection().Send(1, []byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := w2net.NewServer("w2inx v0.5")

	s.AddRouter(&PingRouter{})

	s.Serve()
}
