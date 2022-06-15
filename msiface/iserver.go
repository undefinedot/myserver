package msiface

// IServer 服务器接口，对外开放
type IServer interface {
	Start() // 启动服务器
	Stop()  // 停止服务器
	Serve() // 运行服务器

	AddRouter(msgID uint32, router IRouter) // 路由注册，路由的方法将由连接来使用
}
