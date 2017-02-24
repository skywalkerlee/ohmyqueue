package server

import (
	"golang.org/x/net/context"

	"strconv"

	log "github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/broker"
	"github.com/ohmq/ohmyqueue/clientrpc"
)

type RpcServer struct {
	Broker *broker.Broker
}

func (rpcserver *RpcServer) PutMsg(ctx context.Context, remotemsg *clientrpc.Msg) (*clientrpc.StatusCode, error) {
	log.Info("PutMsg")
	log.Info(remotemsg.GetBody())
	rpcserver.Broker.Put(remotemsg.GetBody())

	return &clientrpc.StatusCode{Code: 200, Offset: strconv.Itoa(rpcserver.Broker.Len())}, nil
}

func (rpcserver *RpcServer) Poll(ctx context.Context, req *clientrpc.Req) (*clientrpc.Resp, error) {
	body, err := rpcserver.Broker.Get(req.GetOffset())
	if err == nil {
		return &clientrpc.Resp{Offset: req.GetOffset(), Msg: body}, err
	}
	return nil, err
}

// func (self *RpcServer) PutMsg(ctx context.Context, remotemsg *clientrpc.Msg) (*clientrpc.StatusCode, error) {
// 	log.Info("PutMsg")
// 	localmsg := msg.Msg{
// 		Header: msg.Header{
// 			Len:      utf8.RuneCountInString(remotemsg.GetBody()),
// 			Deadline: time.Now().Unix() + self.Broker.Msgs.Topics[remotemsg.GetTopic()].Alivetime,
// 		},
// 		Body: remotemsg.GetBody(),
// 	}
// 	log.Info("%s %#v", remotemsg.GetTopic(), localmsg)
// 	self.BrokerPut(remotemsg.GetTopic(), localmsg)
// 	self.Broker.Etcd.Client.Put(context.TODO(), "topic/"+remotemsg.GetTopic()+"/attr", strconv.Itoa(len(self.Broker.Msgs.Topics[remotemsg.GetTopic()].Message)-1))
// 	return &clientrpc.StatusCode{Code: 200}, nil
// }

// func (self *RpcServer) Poll(ctx context.Context, req *clientrpc.Req) (*clientrpc.Resp, error) {
// 	msg := self.Broker.Msgs.Get(req.GetTopic(), req.GetOffset())
// 	return &clientrpc.Resp{
// 		Offset: req.Offset,
// 		Msg: &clientrpc.Msg{
// 			Topic: req.GetTopic(),
// 			Body:  msg.Body,
// 		},
// 	}, nil
// }
