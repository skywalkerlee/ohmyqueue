package main

import (
	"context"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/etcd"
	"github.com/ohmq/ohmyqueue/serverpb"
	"google.golang.org/grpc"
)

func main() {
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	etcd := etcd.NewEtcd()
	resp, _ := etcd.Client.Get(context.TODO(), "topic/test1")
	addr := ""
	for _, ev := range resp.Kvs {
		resp, _ := etcd.Client.Get(context.TODO(), string(ev.Value))
		for _, ev := range resp.Kvs {
			logs.Info(string(ev.Value))
			addr = string(ev.Value)
		}
	}
	conn, _ := grpc.Dial(addr, grpc.WithInsecure())
	client := serverpb.NewOmqClient(conn)
	resp, _ = etcd.Client.Get(context.TODO(), "topic/test1/attr")
	rpcresp := &serverpb.Resp{}
	rpcresp.Offset = "0"
	rpcresp, _ = client.Poll(context.TODO(), &serverpb.Req{Topic: "test1", Offset: rpcresp.Offset})
	logs.Info(rpcresp)
	for _, ev := range resp.Kvs {
		for {
			if rpcresp.GetOffset() == string(ev.Value) {
				break
			}
			offset, _ := strconv.Atoi(rpcresp.GetOffset())
			offset++
			off := strconv.Itoa(offset)
			rpcresp, _ = client.Poll(context.TODO(), &serverpb.Req{Topic: "test1", Offset: off})
			logs.Info(rpcresp)
		}
	}
	wch := etcd.Client.Watch(context.TODO(), "topic/test1/attr")
	for wresp := range wch {
		for _, ev := range wresp.Events {
			switch ev.Type.String() {
			case "PUT":
				offset, _ := strconv.Atoi(rpcresp.GetOffset())
				offset++
				off := strconv.Itoa(offset)
				rpcresp, _ = client.Poll(context.TODO(), &serverpb.Req{Topic: "test1", Offset: off})
				logs.Info(rpcresp)
			}
		}
	}

}
