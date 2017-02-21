package main

import (
	"net"

	"os"

	log "github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/msg"
	"github.com/ohmq/ohmyqueue/server"
	"github.com/ohmq/ohmyqueue/serverpb"
	"google.golang.org/grpc"
)

func main() {
	log.SetLogger("console")
	log.SetLogger(log.AdapterFile, `{"filename":"test.log"}`)
	msgs := msg.NewMsgs()
	broker := server.NewBroker(1, msgs)
	go broker.Start()
	lis, err := net.Listen("tcp", ":9988")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	s := grpc.NewServer()
	serverpb.RegisterOmqServer(s, &server.RpcServer{Msgs: msgs})
	s.Serve(lis)
}
