package routers

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"google.golang.org/protobuf/proto"
	"u3dsrv02/bizcore"
	"u3dsrv02/pb"
)

type SyncPlayerPositionHandle struct {
	znet.BaseRouter
}

func (this *SyncPlayerPositionHandle) Handle(request ziface.IRequest) {
	zlog.Debug("Call SyncPlayerPositionHandle Handle recv from client : msgId=", request.GetMsgID())

	//先读取客户端的数据
	msgRecv := &pb.PlayerInfo{}
	err := proto.Unmarshal(request.GetData(), msgRecv)
	if err != nil {
		zlog.Error(err)
		return
	}

	thisPlayer := bizcore.PlayerMgr.GetPlayerByPID(msgRecv.GetPlayerId())
	if thisPlayer == nil {
		return
	}
	thisPlayer.X = msgRecv.X
	thisPlayer.Y = msgRecv.Y
	thisPlayer.Z = msgRecv.Z
}
