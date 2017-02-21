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
	etcd := etcd.NewEtcd()
	etcd.Client.Put(context.TODO(), "broker/index1/topics", "test1,test2")
	conn, _ := grpc.Dial("127.0.0.1:9988", grpc.WithInsecure())
	client := serverpb.NewOmqClient(conn)
	statuscode, err := client.PutMsg(context.TODO(), &serverpb.Msg{"test1", "this is the first message of omq"})
	if err != nil {
		logs.Error(err)
		os.Exit(1)
	}
	logs.Info(statuscode.GetCode())
}
