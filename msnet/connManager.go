package msnet

import (
	"errors"
	"fmt"
	"myserver/msiface"
	"sync"
)

type ConnManager struct {
	connLock    sync.RWMutex
	connections map[uint32]msiface.IConnection
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]msiface.IConnection),
	}
}

func (cm *ConnManager) Add(conn msiface.IConnection) {
	// 不用defer，因为后面要打印日志才会defer，解锁太慢;其次print调用Count方法中也有锁，死锁
	cm.connLock.Lock()
	cm.connections[conn.GetConnID()] = conn
	cm.connLock.Unlock()

	fmt.Printf("add conn succeed! connID is %v, count = %d\n", conn.GetConnID(), cm.Count())
}

func (cm *ConnManager) Remove(conn msiface.IConnection) {
	cm.connLock.Lock()
	delete(cm.connections, conn.GetConnID())
	cm.connLock.Unlock()

	fmt.Printf("remove conn succeed! connID is %v, count = %d\n", conn.GetConnID(), cm.Count())
}

func (cm *ConnManager) Get(connID uint32) (msiface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connID not exists")
}

func (cm *ConnManager) Count() int {
	// 小心deadlock, Add和Remove中
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	return len(cm.connections)
}

func (cm *ConnManager) CleanConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除map所有conn
	for connID, connection := range cm.connections {
		// 先关闭连接
		connection.Stop()
		delete(cm.connections, connID)
	}
}
