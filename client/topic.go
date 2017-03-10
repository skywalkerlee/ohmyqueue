package main

import (
	"os"

	"time"

	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/clientv3"
	"github.com/ohmq/ohmyqueue/etcd"
	"golang.org/x/net/context"
)

func main() {
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	if len(os.Args) < 2 {
		logs.Error("args err")
		os.Exit(1)
	}
	etcd := etcd.NewEtcd()
	defer etcd.Client.Close()
	if resp, _ := etcd.Client.Txn(context.TODO()).
		If(clientv3.Compare(clientv3.CreateRevision("topicname"+os.Args[1]), "=", 0)).
		Then(clientv3.OpPut("topicname"+os.Args[1], time.Now().String())).Commit(); resp.Succeeded {
		logs.Info("creat topic:", os.Args[1])
	} else {
		logs.Error("topic alread exits")
	}
}
