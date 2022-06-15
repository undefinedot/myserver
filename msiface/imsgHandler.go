package msiface

type IMsgHandler interface {
	DoMsgHandler(request IRequest)          // 根据msgID选择对应的处理函数
	AddRouter(msgID uint32, router IRouter) // 添加路由
}
