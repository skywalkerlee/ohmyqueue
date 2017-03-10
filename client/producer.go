package main

import (
	"context"

	"os"

	"github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/clientrpc"
	"github.com/ohmq/ohmyqueue/etcd"
	"google.golang.org/grpc"
)

func main() {
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	if len(os.Args) < 3 {
		logs.Error("Args error")
		os.Exit(1)
	}
	etcd := etcd.NewEtcd()
	defer etcd.Client.Close()
	resp, _ := etcd.Client.Get(context.TODO(), "topicleader"+os.Args[1])
	logs.Info(string(resp.Kvs[0].Value))
	conn, _ := grpc.Dial(string(resp.Kvs[0].Value), grpc.WithInsecure())
	client := clientrpc.NewOmqClient(conn)
	_, err := client.PutMsg(context.TODO(), &clientrpc.Msg{Topic: os.Args[1], Body: os.Args[2]})
	if err != nil {
		logs.Error(err)
		os.Exit(1)
	}
}
