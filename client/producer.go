package main

import (
	"context"

	"os"

	"github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/etcd"
	"github.com/ohmq/ohmyqueue/serverpb"
	"google.golang.org/grpc"
)

func main() {
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	if len(os.Args) < 2 {
		logs.Error("err")
		os.Exit(1)
	}
	etcd := etcd.NewEtcd()
	etcd.Client.Put(context.TODO(), "broker/index1/topics", "test1")
	etcd.Client.Put(context.TODO(), "topic/test1", "broker/index1")
	conn, _ := grpc.Dial("127.0.0.1:9988", grpc.WithInsecure())
	client := serverpb.NewOmqClient(conn)
	statuscode, err := client.PutMsg(context.TODO(), &serverpb.Msg{"test1", os.Args[1]})
	if err != nil {
		logs.Error(err)
		os.Exit(1)
	}
	logs.Info(statuscode.GetCode())
	resp, _ := etcd.Client.Get(context.TODO(), "topic/test1/attr")
	for _, ev := range resp.Kvs {
		logs.Info(string(ev.Value))
	}
}
