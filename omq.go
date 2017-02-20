package main

import (
	"net"

	"log"

	"github.com/ohmq/ohmyqueue/msg"
	"github.com/ohmq/ohmyqueue/server"
	"github.com/ohmq/ohmyqueue/serverpb"
	"google.golang.org/grpc"
)

func main() {
	msgs := msg.NewMsgs()
	broker := server.NewBroker(1, msgs)
	go broker.Start()
	lis, err := net.Listen("tcp", ":9988")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	serverpb.RegisterOmqServer(s, &server.RpcServer{Msgs: msgs})
	s.Serve(lis)
}
