package main

import (
	"fmt"
	"github.com/lnnlmario/w2-inx/w2net"
	"io"
	"net"
)

// 负责测试datapack拆包，封包功能
func main() {
	//创建socket TCP Server
	listener, err := net.Listen("tcp", "127.0.0.1:8999")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("server accept err:", err)
		}

		go func(conn net.Conn) {
			defer conn.Close()

			// 创建拆包对象
			dp := w2net.NewDataPack()
			for {
				// 读取head部分
				headData := make([]byte, dp.GetHeadLen())
				if _, err := io.ReadFull(conn, headData); err != nil {
					fmt.Println("read head err:", err)
					break
				}
				// 将headData字节流拆包到msg中
				msgHead, err := dp.Unpack(headData)
				if err != nil {
					fmt.Println("unpack err:", err)
					return
				}

				if msgHead.GetDataLen() > 0 {
					//msg 是有data数据的，需要再次读取data数据
					msg := msgHead.(*w2net.Message)
					msg.Data = make([]byte, msg.GetDataLen())

					//根据dataLen从io中读取字节流
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server unpack data err:", err)
						return
					}

					fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
				}
			}
		}(conn)
	}
}
