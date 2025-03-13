package main

import (
	"sync"
)

/*
	当前游戏的世界总管理模块
*/
type WorldManager struct {
	AoiMgr *AOIManager 	        // 当前世界地图的AOI管理模块
	Players map[int32]*Player   // 当前在线的玩家集合
	pLock   sync.RWMutex        //保护Players的读写锁
}

/* 
	初始化方法，创建一个世界总管理模块对象
	对于总管理模块，我们希望只有一个全局唯一对象WorldMgrObj （单例）
	然后我们可以通过init()函数来创建这个对象
*/
var WorldMgrObj *WorldManager
func init() {
	WorldMgrObj = &WorldManager{
		AoiMgr: NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y), // 初始化一个AOI管理模块
		Players: make(map[int32]*Player),
	}
}

// 提供添加一个玩家的的功能，将玩家添加进玩家信息表Players
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()  // 加锁
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()  // 解锁

	// 将player添加到AOIManager中
	wm.AoiMgr.AddToGridByPos(player.Pid, player.X, player.Z)
}

// 提供删除一个玩家的的功能，将玩家从玩家信息表Players删除
func (wm *WorldManager) RemovePlayerByPid(pid int32) {
	// 将player添加到AOIManager中
	wm.pLock.RLock()  // 加读锁
	plyaer := vm.Players[pid]  
	wm.pLock.RUnlock()  // 解读锁
	wm.AoiMgr.RemovePidFromGrid(pid, player.X, player.Z)

	wm.pLock.Lock()  // 加锁
	delete(wm.Players, pid)
	wm.pLock.Unlock()  // 解锁
}

// 通过玩家ID查询Player对象
func (wm *WorldManager) GetPlayerByPid(pid int32) {
	wm.pLock.RLock()  // 加读锁
	defer wm.pLock.RUnlock()  // 解读锁
	return vm.Players[pid]  
}

// 获取全部的在线玩家
func (wm *WorldManager) GetAllPlayer() []*Players {
	wm.pLock.RLock()  // 加读锁
	defer wm.pLock.RUnlock()  // 解读锁
	
	players := make([]*Players, 0)

	// 将所有玩家添加到切片players中
	for _, p := range vm.Players {
		players = append(players, P)
	}

	return players
}