package core

import (
	"fmt"
	"sync"
)

/*
一个AOI地图中的格子数据类型
*/
type Grid struct {
	GID       int          // 格子的ID
	MinX      int          // 格子的左边边界坐标
	MaxX      int          // 格子的右边边界坐标
	MinY      int          // 格子的上边边界坐标
	MaxY      int          // 格子的下边边界坐标
	playerIDs map[int]bool // 格子内玩家/物体成员的ID集合
	pIDLock   sync.RWMutex // playerIDs的保护锁
}

// 创建一个格子
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// 给格子添加一个玩家
func (g *Grid) Add(playerID int) {
	// 加锁保护共享资源
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

// 从格子中删除一个玩家
func (g *Grid) Remove(playerID int) {
	// 加锁保护共享资源
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// 获得当前格子中所有的玩家
func (g *Grid) GetPlyerIDs() (playerIDs []int) {
	// 加锁保护共享资源
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}

	return
}

// 打印出格子的基本信息
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX: %d, maxX: %d, minY: %d, maxY: %d, playerIDs: %v", g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
