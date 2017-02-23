package server

import (
	"golang.org/x/net/context"

	log "github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/broker"
	"github.com/ohmq/ohmyqueue/serverpb"
)

type RpcServer struct {
	Broker *broker.Broker
}

func (rpcserver *RpcServer) PutMsg(ctx context.Context, remotemsg *serverpb.Msg) (*serverpb.StatusCode, error) {
	log.Info("PutMsg")
	err := rpcserver.Broker.Put(remotemsg.GetOffset(), remotemsg.GetBody())
	if err == nil {
		return &serverpb.StatusCode{Code: 200}, err
	}
	return nil, err
}

func (rpcserver *RpcServer) Poll(ctx context.Context, req *serverpb.Req) (*serverpb.Resp, error) {
	body, err := rpcserver.Broker.Get(req.GetOffset())
	if err == nil {
		return &serverpb.Resp{Offset: req.GetOffset(), Msg: body}, err
	}
	return nil, err
}

// func (self *RpcServer) PutMsg(ctx context.Context, remotemsg *serverpb.Msg) (*serverpb.StatusCode, error) {
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
// 	return &serverpb.StatusCode{Code: 200}, nil
// }

// func (self *RpcServer) Poll(ctx context.Context, req *serverpb.Req) (*serverpb.Resp, error) {
// 	msg := self.Broker.Msgs.Get(req.GetTopic(), req.GetOffset())
// 	return &serverpb.Resp{
// 		Offset: req.Offset,
// 		Msg: &serverpb.Msg{
// 			Topic: req.GetTopic(),
// 			Body:  msg.Body,
// 		},
// 	}, nil
// }
