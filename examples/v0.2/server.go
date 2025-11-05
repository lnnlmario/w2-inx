package main

import "github.com/lnnlmario/w2-inx/w2net"

func main() {
	//创建一个server句柄，使用w2inx的api
	s := w2net.NewServer("[w2-inx v0.2]")
	//启动server
	s.Serve()
}
