package msnet

import (
	"fmt"
	"myserver/msiface"
	"strconv"
)

type MsgHandler struct {
	Routers map[uint32]msiface.IRouter // 注册的所有路由
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Routers: make(map[uint32]msiface.IRouter),
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
