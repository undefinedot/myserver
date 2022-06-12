package msiface

import "net"

// IConnection 定义连接的接口
type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn  // 获取原始的socket
	GetConnID() uint32  // 获取当前连接的ID
	RemoteAddr() net.Addr // 远程客户端的Addr
	Send(data []byte) error // 发数据给远程客户端
}

// HandleFunc 定义一个处理连接业务的方法，函数类型
// TODO: 带更改, p6:11:55
type HandleFunc func(*net.TCPConn, []byte, int) error