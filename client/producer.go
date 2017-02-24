package main

import (
	"context"

	"os"

	"github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/clientrpc"
	"github.com/ohmq/ohmyqueue/etcd"
	"google.golang.org/grpc"
)

func test() {
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	if len(os.Args) < 2 {
		logs.Error("err")
		os.Exit(1)
	}
	etcd := etcd.NewEtcd()
	defer etcd.Client.Close()
	resp, _ := etcd.Client.Get(context.TODO(), "leader")
	logs.Info(string(resp.Kvs[0].Value))
	conn, _ := grpc.Dial(string(resp.Kvs[0].Value), grpc.WithInsecure())
	client := clientrpc.NewOmqClient(conn)
	statuscode, err := client.PutMsg(context.TODO(), &clientrpc.Msg{Body: os.Args[1]})
	if err != nil {
		logs.Error(err)
		os.Exit(1)
	}
	logs.Info(statuscode.GetCode())
	etcd.Client.Put(context.TODO(), "topic", statuscode.GetOffset())
}
