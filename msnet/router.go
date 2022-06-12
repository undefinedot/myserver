package msnet

import (
	"myserver/msiface"
)

// BaseRouter 用于给具体业务的Router重写
type BaseRouter struct {}

func (b *BaseRouter) PreHandle(msiface.IRequest) {}

func (b *BaseRouter) Handle(msiface.IRequest) {}

func (b *BaseRouter) PostHandle(msiface.IRequest) {}

