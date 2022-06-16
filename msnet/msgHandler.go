package msnet

import (
	"fmt"
	"myserver/msiface"
	"myserver/utils"
	"strconv"
)

type MsgHandler struct {
	Routers        map[uint32]msiface.IRouter // 注册的所有路由
	WorkerPoolSize uint32                     // 工作池的worker数量
	TaskQueue      []chan msiface.IRequest    // 消息队列
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Routers:        make(map[uint32]msiface.IRouter),
		WorkerPoolSize: utils.GlobalConfig.WorkerPoolSize,
		TaskQueue:      make([]chan msiface.IRequest, utils.GlobalConfig.MaxWorkerTaskLen),
	}
}

func (mh *MsgHandler) DoMsgHandler(request msiface.IRequest) {
	handler, ok := mh.Routers[request.GetMsgID()]
	if !ok {
		fmt.Println("router not found! msgID:", request.GetMsgID())
		return
	}

	// todo: [记]模板方法设计模式，固定执行顺序，用户只能改写方法的具体实现
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router msiface.IRouter) {
	if _, ok := mh.Routers[msgID]; ok {
		panic("message id exists, id = " + strconv.Itoa(int(msgID)))
	}
	// 不存在才添加
	mh.Routers[msgID] = router
	fmt.Println("add router succeed, msgID:", msgID)
}

// StartWorkerPool 开启工作池
func (mh *MsgHandler) StartWorkerPool() {
	// 只能开启一个工作池
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		fmt.Printf("开启第 %v 个worker...\n", i)
		// 每个worker对应一个消息队列, i作为它的ID
		mh.TaskQueue[i] = make(chan msiface.IRequest, utils.GlobalConfig.MaxWorkerTaskLen)
		go mh.StartOneWorker(i)
	}
}

// StartOneWorker 开启一个worker
func (mh *MsgHandler) StartOneWorker(workerID int) {
	for {
		select {
		case req := <-mh.TaskQueue[workerID]:
			mh.DoMsgHandler(req)
		}
	}
}

// SendReqToTaskQueue 将request发送给task Queue
func (mh *MsgHandler) SendReqToTaskQueue(req msiface.IRequest) {
	// 保证分配均衡，// todo: 轮询的平均分配法则, 可优化:根据req的ID
	workerID := req.GetConnection().GetConnID() % utils.GlobalConfig.WorkerPoolSize
	mh.TaskQueue[workerID] <- req
}
