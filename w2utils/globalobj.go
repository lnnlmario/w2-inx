package w2utils

import (
	"encoding/json"
	"github.com/lnnlmario/w2-inx/w2iface"
	"io/ioutil"
)

// 存储矿建全局参数
type GlobalObj struct {
	TcpServer     w2iface.IServer // 全局Server对象
	Host          string          // 服务主机IP
	TcpPort       int             // 监听端口
	Name          string          // 服务器名称
	Version       string          // w2inx 版本号
	MaxPacketSize uint32          // 数据包最大值
	MaxConn       int             // 允许的最大链接数
}

var GlobalObject *GlobalObj

// 读取用户的配置文件
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/w2inx.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	//fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:          "W2inxServerApp",
		Version:       "v0.4",
		TcpPort:       8999,
		Host:          "0.0.0.0",
		MaxPacketSize: 4096,
		MaxConn:       12000,
	}

	//从配置文件中加载一些用户配置的参数
	GlobalObject.Reload()
}
