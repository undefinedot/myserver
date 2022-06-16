package msiface

type IConnManager interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	Get(connID uint32) (IConnection, error)
	Count() int // 计算当前连接总数
	CleanConn() // 清除所有连接
}
