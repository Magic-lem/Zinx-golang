package znet

import (
	"fmt"
	"sync"
	"errors"
	"workspace/src/zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection // 保存ID-连接的Map
	connLock	sync.RWMutex				  // 保护连接Map的读写锁
}

// 构造函数：创建一个ConnManager实例
func NewConnManager() *ConnManager {
	return &ConnManager {
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加连接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源Map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()  // 函数结束时解锁

	cm.connections[conn.GetConnID()] = conn
	fmt.Println("ConnID = ", conn.GetConnID(), "add to ConnManager successfully, conn num = ", cm.Len())
}

// 删除连接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源Map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()  // 函数结束时解锁

	delete(cm.connections, conn.GetConnID())
	fmt.Println("ConnID = ", conn.GetConnID(), "remove from ConnManager successfully, conn num = ", cm.Len())
}

// 根据连接ID获取连接
func (cm *ConnManager) GetConn(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源Map，加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()  // 函数结束时解锁

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil,  errors.New("connection not Found!")
	}
}

// 当前连接个数
func (cm *ConnManager) Len() int {
	// 不需要加锁，否则会导致死锁
	// TODO：不加锁的话则需要保证Len()函数只会在加锁后的环境中调用，如Add、Remove等方法中，保证并发安全
	return len(cm.connections)
}

// 清除所有连接
func (cm *ConnManager) ClearConn() {
	// 保护共享资源Map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()  // 函数结束时解锁

	// 删除所有conn并停止其工作
	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}
	
	fmt.Println("Clear All connections succ! conn num = ", cm.Len())
}	