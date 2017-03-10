package server

import (
	"golang.org/x/net/context"

	"time"

	"strconv"

	log "github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/broker"
	"github.com/ohmq/ohmyqueue/clientrpc"
	"github.com/ohmq/ohmyqueue/config"
)

type RpcServer struct {
	Broker *broker.Broker
}

func (rpcserver *RpcServer) PutMsg(ctx context.Context, remotemsg *clientrpc.Msg) (*clientrpc.StatusCode, error) {
	log.Info("PutMsg")
	log.Info("alivetime", config.Conf.Topic.Alivetime)
	rpcserver.Broker.Put(remotemsg.GetTopic(), strconv.FormatInt(time.Now().Unix()+config.Conf.Topic.Alivetime, 10), remotemsg.GetBody())
	return &clientrpc.StatusCode{Code: 200}, nil
}

func (rpcserver *RpcServer) Poll(ctx context.Context, req *clientrpc.Req) (*clientrpc.Resp, error) {
	log.Info("Poll")
	offset, body, err := rpcserver.Broker.Get(req.GetTopic(), req.GetOffset())
	return &clientrpc.Resp{Body: body, Offset: offset}, err
}
