package w2net

import (
	"fmt"

	"github.com/lnnlmario/w2-inx/w2iface"
)

type MsgHandle struct {
	// 存放每个msgId对应的处理方法映射
	Apis map[uint32]w2iface.IRouter
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]w2iface.IRouter),
	}
}

func (m *MsgHandle) DoMsgHandler(request w2iface.IRequest) {
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId=", request.GetMsgID(), " not exist!")
		return
	}

	// 执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (m *MsgHandle) AddRouter(msgId uint32, router w2iface.IRouter) {
	// 绑定api是否存在
	if _, ok := m.Apis[msgId]; ok {
		panic("repeated api, msgId=" + fmt.Sprint(msgId))
	}

	m.Apis[msgId] = router
	fmt.Println("Add Router msgId= ", msgId)
}
