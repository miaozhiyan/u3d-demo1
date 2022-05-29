package bizcore

import "github.com/aceld/zinx/ziface"

type Player struct {
	PlayerId uint32
	X        float32
	Y        float32
	Z        float32
	Conn ziface.IConnection //当前玩家的连接
}
