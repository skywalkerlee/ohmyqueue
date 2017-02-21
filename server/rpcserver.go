package server

import (
	"time"
	"unicode/utf8"

	log "github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/msg"
	"github.com/ohmq/ohmyqueue/serverpb"
	"golang.org/x/net/context"
)

type RpcServer struct {
	Msgs *msg.Msgs
}

func (self *RpcServer) PutMsg(ctx context.Context, remotemsg *serverpb.Msg) (*serverpb.StatusCode, error) {
	log.Info("PutMsg")
	localmsg := msg.Msg{
		Header: msg.Header{
			Len:      utf8.RuneCountInString(remotemsg.GetBody()),
			Deadline: time.Now().Unix() + self.Msgs.Topics[remotemsg.GetTopic()].Alivetime,
		},
		Body: remotemsg.GetBody(),
	}
	log.Info("%s %#v", remotemsg.GetTopic(), localmsg)
	self.Msgs.Put(remotemsg.GetTopic(), localmsg)
	return &serverpb.StatusCode{Code: 200}, nil
}

func (self *RpcServer) Poll(ctx context.Context, req *serverpb.Req) (*serverpb.Resp, error) {
	msg := self.Msgs.Get(req.GetTopic(), req.GetOffset())
	return &serverpb.Resp{
		Offset: req.Offset,
		Msg: &serverpb.Msg{
			Topic: req.GetTopic(),
			Body:  msg.Body,
		},
	}, nil
}
