package msiface

import "net"

// IConnection 定义连接的接口
type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn                      // 获取原始的socket
	GetConnID() uint32                                   // 获取当前连接的ID
	RemoteAddr() net.Addr                                // 远程客户端的Addr
	SendMsg(msgID uint32, data []byte) error             // 发数据给远程客户端
	SetProperty(key string, val interface{})             // 设置连接属性
	GetProperty(key string) (val interface{}, err error) // 获取连接属性
	RemoveProperty(key string)                           // 移除连接属性
}
