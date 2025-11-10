package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Client start...")

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		_, err := conn.Write([]byte("hello world"))
		if err != nil {
			fmt.Println(err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("server call back: %s, cnt: %d\n", buf, cnt)
		time.Sleep(1 * time.Second)
	}
}
