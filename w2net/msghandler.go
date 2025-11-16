package w2net

import (
	"fmt"
	"github.com/lnnlmario/w2-inx/w2utils"

	"github.com/lnnlmario/w2-inx/w2iface"
)

type MsgHandle struct {
	// 存放每个msgId对应的处理方法映射
	Apis           map[uint32]w2iface.IRouter
	WorkerPoolSize uint32                  // 业务工作Worker池的数量
	TaskQueue      []chan w2iface.IRequest // Worker 负责取任务的消息队列
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]w2iface.IRouter),
		WorkerPoolSize: w2utils.GlobalObject.WorkerPoolSize,
		// 缓冲worker调用的Request请求信箱，worker回从对应的队列中获取数据
		TaskQueue: make([]chan w2iface.IRequest, w2utils.GlobalObject.WorkerPoolSize),
	}
}

func (m *MsgHandle) DoMsgHandler(request w2iface.IRequest) {
	// 通过消息id得到对应处理方法
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

func (m *MsgHandle) StartOneWorker(workerId int, taskQueue chan w2iface.IRequest) {
	fmt.Println("StartOneWorker workerId=", workerId)
	// 不断等待队列中的消息
	for {
		select {
		// 有消息则取出队列中的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// 启动worker工作池
func (m *MsgHandle) StartWorkerPool() {
	//遍历需要启动worker的数量，依次启动
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 启动一个worker
		// 给当前worker对应的任务队列开辟空间
		m.TaskQueue[i] = make(chan w2iface.IRequest, w2utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前worker 阻塞瞪大对应的任务队列是否有消息传递进来
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// 将消息交给TaskQuque,由worker进行处理
func (m *MsgHandle) SendMsgToTaskQueue(req w2iface.IRequest) {
	// 根据connId来分配当前的链接应该由哪个worker负责处理

	// 轮询的平均分配法则
	// 得到需要处理链接的workerID
	workerId := req.GetConnection().GetConnId() % m.WorkerPoolSize
	fmt.Println("SendMsgToTaskQueue workerId=", workerId, "Request msgId=", req.GetMsgID())
	// 将请求消息发送给任务队列
	m.TaskQueue[workerId] <- req
}
