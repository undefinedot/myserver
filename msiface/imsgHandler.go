package msiface

type IMsgHandler interface {
	DoMsgHandler(request IRequest)          // 根据msgID选择对应的处理函数
	AddRouter(msgID uint32, router IRouter) // 添加路由
	StartWorkerPool()                       // 开启工作池
	SendReqToTaskQueue(request IRequest)    // 发送request给消息队列
}
