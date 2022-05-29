package routers

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"google.golang.org/protobuf/proto"
	"u3dsrv02/bizcore"
	"u3dsrv02/pb"
)

type InitPositionRouter struct {
	znet.BaseRouter
}

//HelloZinxRouter Handle
func (this *InitPositionRouter) Handle(request ziface.IRequest) {
	zlog.Debug("Call InitPositionRouter Handle recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	//先读取客户端的数据
	msgRecv := &pb.PlayerInfo{}
	err := proto.Unmarshal(request.GetData(), msgRecv)
	if err != nil {
		zlog.Error(err)
		return
	}

	// 向所有其他玩家同步这个新的玩家的加入
	allPlayers := bizcore.PlayerMgr.GetAllPlayers()
	for _, p := range allPlayers {
		if p.PlayerId == msgRecv.GetPlayerId() {
			// 自己的位置信息没必要同步
			continue
		}
		err = p.Conn.SendBuffMsg(50, request.GetData())
		if err != nil {
			zlog.Error(err)
		}
	}

	// 向这个新的玩家同步其他之前已经存在的玩家的信息
	playersData := make([]*pb.PlayerInfo, 0, len(allPlayers))
	for _, playInfoInServer := range allPlayers {
		p := &pb.PlayerInfo{
			X:        playInfoInServer.X,
			Y:        playInfoInServer.Y,
			Z:        playInfoInServer.Z,
			PlayerId: playInfoInServer.PlayerId,
		}
		playersData = append(playersData, p)
	}

	synOthers := &pb.SynOtherPlayerInfos{
		PlayerInfos: playersData[:],
	}
	synOthersMsg, err := proto.Marshal(synOthers)
	if err != nil {
		zlog.Error(err)
		return
	}
	err = request.GetConnection().SendBuffMsg(51, synOthersMsg)
	if err != nil {
		zlog.Error(err)
	}
}
