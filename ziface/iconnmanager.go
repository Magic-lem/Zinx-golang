package ziface

type IConnManager interface {
	// 添加连接
	Add(conn IConnection)
	// 删除连接
	Remove(conn IConnection)
	// 根据连接ID获取连接
	GetConn(connID uint32) (IConnection, error)
	// 当前连接个数
	Len() int
	// 清除所有连接
	ClearConn()
}