package core

import (
	"fmt"
)

/*
	AOI管理模块
*/
type AOIManager struct {
	MinX		int    			// 区域的左边界坐标​
	MaxX		int				// 区域的右边界坐标​
	CntsX		int				// X方向格子的数量
	MinY		int				// 区域的上边界坐标​
	MaxY		int				// 区域的下边界坐标​
	CntsY		int				// Y方向格子的数量
	grids       map[int] *Grid	// 当前区域中都有哪些格子，key=格子ID， value=格子对象
}

// 初始化一个AOI区域
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX: minX,
		MaxX: maxX,
		CntsX: cntsX,
		MinY: minY,
		MaxY: maxX,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	// 为AOI初始化区域中的所有格子
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// 计算格子的ID：id = idy * nx + idx
			gid := y * cntsX + x

			// 初始化一个格子放在AOI的map中
			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX + x * aoiMgr.gridWidth(),
				aoiMgr.MinX + (x + 1) * aoiMgr.gridWidth(),
				aoiMgr.MinY + y * aoiMgr.gridLength(),
				aoiMgr.MinY + (y + 1) * aoiMgr.gridLength(),
				)
		}
	}

	return aoiMgr
}

// 得到每个格子在x轴方向的宽度
func (aoiMgr *AOIManager) gridWidth() int {
	return (aoiMgr.MaxX - aoiMgr.MinX) / aoiMgr.CntsX
}

//得到每个格子在x轴方向的长度
func (aoiMgr *AOIManager) gridLength() int {
	return (aoiMgr.MaxY - aoiMgr.MinY) / aoiMgr.CntsY
}

// 调试使用：打印当前AOI模块
func (aoiMgr *AOIManager) String() string {
	s := fmt.Sprintf("AOIManagr:\nminX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d\n Grids in AOI Manager:\n", 
					aoiMgr.MinX, aoiMgr.MaxX, aoiMgr.CntsX, aoiMgr.MinY, aoiMgr.MaxY, aoiMgr.CntsY)
	for _, grid := range aoiMgr.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

// 根据格子ID求出其感兴趣的九宫格
func (aoiMgr *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	// 判断输入的gID是否在AOIManager中
	if _, ok := aoiMgr.grids[gID]; !ok {
		return
	}

	// 初始化grids，将当前gid添加到九宫格中
	grids = append(grids, aoiMgr.grids[gID])

	// 判断gID的左边是否有格子，右边是否有格子？
	idx := gID % aoiMgr.CntsX	// 利用ID求出x轴编号
	if idx > 0 {
		grids = append(grids, aoiMgr.grids[gID - 1])
	}
	if idx < aoiMgr.CntsX - 1 {
		grids = append(grids, aoiMgr.grids[gID + 1])
	}

	// 将x轴当前的格子都取出，进行遍历，再分别得到每个格子的上下是否有格子
	gridsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gridsX = append(gridsX, v.GID)
	}

	// 遍历X轴上的所有格子
	for _, v := range gridsX {
		// 计算出格子所在的列
		idy := v / aoiMgr.CntsX

		if idy > 0 {
			grids = append(grids, aoiMgr.grids[v - aoiMgr.CntsX])
		}
		if idy < aoiMgr.CntsY - 1 {
			grids = append(grids, aoiMgr.grids[v + aoiMgr.CntsX])
		}
	}

	return
}

// 通过横纵坐标获取对应的格子ID
func (aoiMgr *AOIManager) GetGIDByPos(x, y float32) int {
	gx := (int(x) - aoiMgr.MinX) / aoiMgr.gridWidth()
	gy := (int(y) - aoiMgr.MinY) / aoiMgr.gridLength()

	return gy * aoiMgr.CntsX + gx
}

// 通过横纵坐标得到周边九宫格内的全部PlayerIDs
func (aoiMgr *AOIManager) GetPIDsByPos(x, y float32) (playerIDs []int) {
	// 求出坐标所在的格子ID
	gID := aoiMgr.GetGIDByPos(x, y)

	// 根据格子ID获取周边九宫格
	grids := aoiMgr.GetSurroundGridsByGid(gID)

	// 遍历九宫格来搜集所有玩家
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlyerIDs()...)
		fmt.Printf("===> grid ID : %d, pids : %v  ====", grid.GID, grid.GetPlyerIDs())
	}

	return
}

// 通过GID获取当前格子的全部playerID
func (aoiMgr *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	if grid, ok := aoiMgr.grids[gID]; ok {
		playerIDs = grid.GetPlyerIDs()
	}
	return
}

// 移除一个格子中的PlayerID
func (aoiMgr *AOIManager) RemovePidFromGrid(pID, gID int) {
	if _, ok := aoiMgr.grids[gID]; ok {
		aoiMgr.grids[gID].Remove(pID)
	}
}

// 添加一个PlayerID到一个格子中
func (aoiMgr *AOIManager) AddPidToGrid(pID, gID int) {
	if _, ok := aoiMgr.grids[gID]; ok {
		aoiMgr.grids[gID].Add(pID)
	}
}

// 通过横纵坐标添加一个Player到一个格子中
func (aoiMgr *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := aoiMgr.GetGIDByPos(x, y)
	aoiMgr.grids[gID].Add(pID)
}

// 通过横纵坐标把一个Player从对应的格子中删除
func (aoiMgr *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := aoiMgr.GetGIDByPos(x, y)
	aoiMgr.grids[gID].Remove(pID)
}