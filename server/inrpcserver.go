package server

import (
	"io"

	"github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/broker"
	"github.com/ohmq/ohmyqueue/inrpc"
	"github.com/ohmq/ohmyqueue/msg"
)

type InrpcServer struct {
	topics msg.Topics
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
		inserver.topics.Put(msg.GetTopic(), msg.GetAlivetime(), msg.GetBody(), msg.GetOffset())
	}
}
