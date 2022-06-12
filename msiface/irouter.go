package msiface

type IRouter interface {
	PreHandle(r IRequest) // 前钩子方法
	Handle(r IRequest) // 处理conn业务的方法
	PostHandle(r IRequest) // 后狗子方法
}
