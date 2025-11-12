package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/lnnlmario/w2-inx/w2net"
)

func main() {
	// 3s后发起测试请求，给服务器开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start error:", err)
		return
	}

	for {
		// 发送messge消息
		dp := w2net.NewDataPack()
		msg, _ := dp.Pack(w2net.NewMsgPackage(0, []byte("w2inx v0.5 client test message")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("client write error:", err)
			return
		}

		// 读取流中head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) // ReadFull 会把msg填充满为止
		if err != nil {
			fmt.Println("client write error:", err)
			break
		}

		// 将headData字节流拆到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("client unpack error:", err)
			return
		}

		// 判断是否有data数据
		if msgHead.GetDataLen() > 0 {
			msg := msgHead.(*w2net.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			// 根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("client read error:", err)
				return
			}

			fmt.Println("==> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}
