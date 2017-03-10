package main

import (
	"context"
	"os"

	"strings"

	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/clientrpc"
	"github.com/ohmq/ohmyqueue/etcd"
	"google.golang.org/grpc"
)

//main
func main() {
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	if len(os.Args) < 2 {
		logs.Error("Args error")
		os.Exit(1)
	}
	etcd := etcd.NewEtcd()
	defer etcd.Client.Close()
	resp, _ := etcd.Client.Get(context.TODO(), "topicleader"+os.Args[1])
	logs.Info(string(resp.Kvs[0].Value))
	conn, _ := grpc.Dial(string(resp.Kvs[0].Value), grpc.WithInsecure())
	client := clientrpc.NewOmqClient(conn)
	resp, _ = etcd.Client.Get(context.TODO(), "topicattr"+os.Args[1])
	if resp.Count != 0 {
		offmax, _ := strconv.Atoi(strings.Split(string(resp.Kvs[0].Value), ":")[0])
		off := 0
		for {
			resp, err := client.Poll(context.TODO(), &clientrpc.Req{Topic: os.Args[1], Offset: int64(off)})
			if err != nil {
				logs.Error(err)
				break
			}
			logs.Info(resp.GetOffset(), resp.GetBody())
			if resp.GetOffset() == -1 || off == offmax {
				break
			}
			off = int(resp.GetOffset()) + 1
		}
	}
	wch := etcd.Client.Watch(context.TODO(), "topicattr"+os.Args[1])
	for wresp := range wch {
		for _, ev := range wresp.Events {
			switch ev.Type.String() {
			case "PUT":
				offmax, _ := strconv.ParseInt((strings.Split(string(ev.Kv.Value), ":")[0]), 10, 64)
				resp, _ := client.Poll(context.TODO(), &clientrpc.Req{Topic: os.Args[1], Offset: offmax})
				logs.Info(resp.GetOffset(), resp.GetBody())
			}
		}
	}
}
