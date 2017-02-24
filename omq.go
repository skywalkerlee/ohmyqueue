package main

import (
	"net"

	"os"

	"strconv"

	log "github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/broker"
	"github.com/ohmq/ohmyqueue/clientrpc"
	"github.com/ohmq/ohmyqueue/inrpc"
	"github.com/ohmq/ohmyqueue/server"
	"google.golang.org/grpc"
)

func main() {
	log.SetLogger("console")
	log.SetLogger(log.AdapterFile, `{"filename":"test.log"}`)
	log.EnableFuncCallDepth(true)
	log.SetLogFuncCallDepth(3)
	if len(os.Args) < 4 {
		log.Error("err")
		os.Exit(1)
	}
	index, _ := strconv.Atoi(os.Args[1])
	cliport, _ := strconv.Atoi(os.Args[2])
	inport, _ := strconv.Atoi(os.Args[3])
	broker := broker.NewBroker(index, cliport, inport)
	go broker.Start()
	lis, err := net.Listen("tcp", ":"+os.Args[2])
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	s := grpc.NewServer()
	clientrpc.RegisterOmqServer(s, &server.RpcServer{Broker: broker})
	go s.Serve(lis)
	lis2, err := net.Listen("tcp", ":"+os.Args[3])
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	s2 := grpc.NewServer()
	inrpc.RegisterInServer(s2, &server.InrpcServer{Broker: broker})
	s2.Serve(lis2)
}
