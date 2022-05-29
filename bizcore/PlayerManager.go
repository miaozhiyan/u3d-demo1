package bizcore

import (
	"sync"
)

type PlayerManager struct {
	Players map[uint32]*Player //当前在线的玩家集合
	pLock   sync.RWMutex       //保护Players的互斥读写机制
}

//提供一个对外的playerManager管理模块句柄
var PlayerMgr *PlayerManager

func init() {
	PlayerMgr = &PlayerManager{
		Players: make(map[uint32]*Player),
	}
}

//提供添加一个玩家的的功能，将玩家添加进玩家信息表Players
func (pm *PlayerManager) AddPlayer(player *Player) {
	//将player添加到 世界管理器中
	pm.pLock.Lock()
	pm.Players[player.PlayerId] = player
	pm.pLock.Unlock()
}

//从玩家信息表中移除一个玩家
func (pm *PlayerManager) RemovePlayerByPID(playerId uint32) {
	pm.pLock.Lock()
	delete(pm.Players, playerId)
	pm.pLock.Unlock()
}

//通过玩家ID 获取对应玩家信息
func (pm *PlayerManager) GetPlayerByPID(playerId uint32) *Player {
	pm.pLock.RLock()
	defer pm.pLock.RUnlock()

	return pm.Players[playerId]
}

//获取所有玩家的信息
func (pm *PlayerManager) GetAllPlayers() []*Player {
	pm.pLock.RLock()
	defer pm.pLock.RUnlock()

	//创建返回的player集合切片
	players := make([]*Player, 0)

	//添加切片
	for _, v := range pm.Players {
		players = append(players, v)
	}

	//返回
	return players
}
