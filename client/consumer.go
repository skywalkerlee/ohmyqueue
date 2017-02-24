package main

import (
	"context"

	"sync"

	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/clientv3"
	"github.com/ohmq/ohmyqueue/clientrpc"
	"github.com/ohmq/ohmyqueue/etcd"
	"google.golang.org/grpc"
)

type cli struct {
	offset  string
	mutex   *sync.Mutex
	etcdcli *clientv3.Client
	client  clientrpc.OmqClient
}

func newcli() *cli {
	return &cli{
		offset:  "0",
		mutex:   new(sync.Mutex),
		etcdcli: etcd.NewEtcd().Client,
	}
}

func (cli *cli) watchleader() {
	wch := cli.etcdcli.Watch(context.TODO(), "leader")
	for wresp := range wch {
		for _, ev := range wresp.Events {
			switch ev.Type.String() {
			case "PUT":
				conn, _ := grpc.Dial(string(ev.Kv.Value), grpc.WithInsecure())
				cli.mutex.Lock()
				cli.client = clientrpc.NewOmqClient(conn)
				cli.mutex.Unlock()
			}
		}
	}
}

func (cli *cli) poll(offset string) *clientrpc.Resp {
	cli.mutex.Lock()
	resp, _ := cli.client.Poll(context.TODO(), &clientrpc.Req{Offset: offset})
	cli.mutex.Unlock()
	return resp
}

func main() {
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	cli := newcli()
	getresp, _ := cli.etcdcli.Get(context.TODO(), "leader")
	conn, _ := grpc.Dial(string(getresp.Kvs[0].Value), grpc.WithInsecure())
	cli.client = clientrpc.NewOmqClient(conn)
	go cli.watchleader()
	resp, _ := cli.etcdcli.Get(context.TODO(), "topic")
	offset := string(resp.Kvs[0].Value)
	if offset != "0" {
		for cli.offset < offset {
			plresp, _ := cli.client.Poll(context.TODO(), &clientrpc.Req{Offset: cli.offset})
			logs.Info("%#v", plresp)
			tmp, _ := strconv.Atoi(cli.offset)
			tmp++
			cli.offset = strconv.Itoa(tmp)
		}
	}
	wch := cli.etcdcli.Watch(context.TODO(), "topic")
	for wresp := range wch {
		for _, ev := range wresp.Events {
			switch ev.Type.String() {
			case "PUT":
				tmp, _ := strconv.Atoi(string(ev.Kv.Value))
				cli.offset = strconv.Itoa(tmp - 1)
				plresp, _ := cli.client.Poll(context.TODO(), &clientrpc.Req{Offset: cli.offset})
				logs.Info("%#v", plresp)
			}
		}
	}
}
