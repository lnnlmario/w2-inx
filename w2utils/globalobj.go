package w2utils

import (
	"encoding/json"
	"fmt"
	"github.com/lnnlmario/w2-inx/w2iface"
	"io/ioutil"
)

// 存储矿建全局参数
type GlobalObj struct {
	//Server
	TcpServer w2iface.IServer // 全局Server对象
	Host      string          // 服务主机IP
	TcpPort   int             // 监听端口
	Name      string          // 服务器名称

	//w2inx
	Version          string // w2inx 版本号
	MaxPacketSize    uint32 // 数据包最大值
	MaxConn          int    // 允许的最大链接数
	WorkerPoolSize   uint32 // 业务工作Worker池的数量
	MaxWorkerTaskLen uint32 // 业务工作Worker对应负责的任务队列最大任务存储数量

	ConfFilePath string
}

var GlobalObject *GlobalObj

// 读取用户的配置文件
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		fmt.Println("can not read ", g.ConfFilePath)
	} else {
		//将json数据解析到struct中
		//fmt.Printf("json :%s\n", data)
		err = json.Unmarshal(data, &GlobalObject)
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:             "W2inxServerApp",
		Version:          "v0.8",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxPacketSize:    4096,
		MaxConn:          12000,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
		ConfFilePath:     "conf/w2inx.json",
	}

	//从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}
