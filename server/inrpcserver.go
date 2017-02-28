package server

import (
	"io"

	"golang.org/x/net/context"

	"github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/broker"
	"github.com/ohmq/ohmyqueue/inrpc"
)

type InrpcServer struct {
	Broker *broker.Broker
}

func (inserver *InrpcServer) SyncMsg(steam inrpc.In_SyncMsgServer) error {
	logs.Info("recv from leader")
	var sum int32
	for {
		msg, err := steam.Recv()
		if err == io.EOF {
			return steam.SendAndClose(&inrpc.StatusCode{Sum: sum})
		}
		if err != nil {
			return err
		}
		sum++
		logs.Info(msg.GetBody())
		inserver.Broker.Put(msg.GetBody())
	}
}

func (inserver *InrpcServer) Del(ctx context.Context, delmsg *inrpc.DelMsg) (*inrpc.StatusCode, error) {
	inserver.Broker.Del(delmsg.GetOffset())
	return &inrpc.StatusCode{Sum: 1}, nil
}
