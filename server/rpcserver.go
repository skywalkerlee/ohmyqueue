package server

import (
	"time"
	"unicode/utf8"

	"strconv"

	log "github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/msg"
	"github.com/ohmq/ohmyqueue/serverpb"
	"golang.org/x/net/context"
)

type RpcServer struct {
	Broker *Broker
}

func (self *RpcServer) PutMsg(ctx context.Context, remotemsg *serverpb.Msg) (*serverpb.StatusCode, error) {
	log.Info("PutMsg")
	localmsg := msg.Msg{
		Header: msg.Header{
			Len:      utf8.RuneCountInString(remotemsg.GetBody()),
			Deadline: time.Now().Unix() + self.Broker.msgs.Topics[remotemsg.GetTopic()].Alivetime,
		},
		Body: remotemsg.GetBody(),
	}
	log.Info("%s %#v", remotemsg.GetTopic(), localmsg)
	self.Broker.msgs.Put(remotemsg.GetTopic(), localmsg)
	self.Broker.etcd.Client.Put(context.TODO(), "topic/"+remotemsg.GetTopic()+"/attr", strconv.Itoa(len(self.Broker.msgs.Topics[remotemsg.GetTopic()].Message)-1))
	return &serverpb.StatusCode{Code: 200}, nil
}

func (self *RpcServer) Poll(ctx context.Context, req *serverpb.Req) (*serverpb.Resp, error) {
	msg := self.Broker.msgs.Get(req.GetTopic(), req.GetOffset())
	return &serverpb.Resp{
		Offset: req.Offset,
		Msg: &serverpb.Msg{
			Topic: req.GetTopic(),
			Body:  msg.Body,
		},
	}, nil
}
