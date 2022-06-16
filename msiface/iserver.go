package msiface

// IServer 服务器接口，对外开放
type IServer interface {
	Start() // 启动服务器
	Stop()  // 停止服务器
	Serve() // 运行服务器

	AddRouter(msgID uint32, router IRouter) // 路由注册，路由的方法将由连接来使用
	GetConnMgr() IConnManager               // 获取连接管理器

	SetOnConnStart(func(IConnection)) //设置该Server的连接创建时Hook函数
	SetOnConnStop(func(IConnection))  //设置该Server的连接断开时的Hook函数
	CallOnConnStart(conn IConnection) //调用连接OnConnStart Hook函数
	CallOnConnStop(conn IConnection)  //调用连接OnConnStop Hook函数
}
