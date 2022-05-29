/**
* @Author: Aceld
* @Date: 2020/12/24 00:24
* @Mail: danbing.at@gmail.com
*    zinx server demo
 */
package main

import (
	"github.com/aceld/zinx/examples/zinx_server/zrouter"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"google.golang.org/protobuf/proto"
	"time"
	"u3dsrv02/bizcore"
	"u3dsrv02/pb"
	"u3dsrv02/routers"
)

//创建连接的时候执行
func DoConnectionBegin(conn ziface.IConnection) {
	zlog.Debug("DoConnecionBegin is Called ... connId", conn.GetConnID())

	//设置链接属性，在连接创建之后
	conn.SetProperty("createTime", time.Now().Unix())

	playerInfo := &pb.PlayerInfo{
		PlayerId: conn.GetConnID(),
		X:        1,
	}
	player := bizcore.Player{
		X:        playerInfo.X,
		Y:        playerInfo.Y,
		Z:        playerInfo.Z,
		PlayerId: playerInfo.PlayerId,
		Conn:     conn,
	}
	bizcore.PlayerMgr.AddPlayer(&player)
	msgSend, err := proto.Marshal(playerInfo)
	// 返回connId
	err = conn.SendMsg(2, msgSend)
	if err != nil {
		zlog.Error(err)
	}
}

//连接断开的时候执行
func DoConnectionLost(conn ziface.IConnection) {

	bizcore.PlayerMgr.RemovePlayerByPID(conn.GetConnID())
	//在连接销毁之前，查询conn的Name，Home属性
	zlog.Debug("DoConneciotnLost is Called ... ,connId", conn.GetConnID())
}

func main() {
	//创建一个server句柄
	s := znet.NewServer()

	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//配置路由
	s.AddRouter(0, &zrouter.PingRouter{})
	s.AddRouter(1, &zrouter.HelloZinxRouter{})

	// 玩家移动信息处理
	s.AddRouter(2, &routers.MoveRouter{})
	// 玩家初始位置信息处理
	s.AddRouter(3, &routers.InitPositionRouter{})
	// 维护各个玩家的位置信息
	s.AddRouter(4, &routers.SyncPlayerPositionHandle{})
	//开启服务
	s.Serve()
}
