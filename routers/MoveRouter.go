package routers

import (
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
	"google.golang.org/protobuf/proto"
	"u3dsrv02/bizcore"
	"u3dsrv02/pb"
)

type MoveRouter struct {
	znet.BaseRouter
}

func (this *MoveRouter) Handle(request ziface.IRequest) {
	zlog.Debug("Call MoveRouter Handle recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	//先读取客户端的数据
	msgRecv := &pb.MoveApi{}
	err := proto.Unmarshal(request.GetData(), msgRecv)

	// 广播客户端
	broadMsg := &pb.BroadCastMove{
		Direct:   msgRecv.Direct,
		PlayerId: request.GetConnection().GetConnID(),
		V:        msgRecv.V,
	}
	msgSend, err := proto.Marshal(broadMsg)
	// todo fmt.Println("测试接收到数据:=======", msgRecv.Direct)
	if err != nil {
		zlog.Error(err)
		return
	}

	allPlayers := bizcore.PlayerMgr.GetAllPlayers()
	for _, p := range allPlayers {
		err = p.Conn.SendBuffMsg(4, msgSend)
		if err != nil {
			zlog.Error(err)
		}
	}

}
