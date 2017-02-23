package broker

import (
	"context"
	"math/rand"
	"os"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/clientv3"

	"strconv"
)

func (broker *Broker) Start() {
	go broker.Etcd.Heartbeat("broker"+strconv.Itoa(broker.id), broker.ip+":"+strconv.Itoa(broker.innerport), 10)
	broker.watchLeader()
}

func (broker *Broker) watchLeader() {
	resp, _ := broker.Etcd.Client.Get(context.TODO(), "leader")
	if resp.Count == 0 {
		go broker.vote()
	} else {
		logs.Info("leader is:", string(resp.Kvs[0].Value))
		broker.leader = string(resp.Kvs[0].Value)
	}
	wch := broker.Etcd.Client.Watch(context.TODO(), "leader")
	for wresp := range wch {
		for _, ev := range wresp.Events {
			switch ev.Type.String() {
			case "PUT":
				logs.Info("leader is:", string(ev.Kv.Value))
				if broker.leader != string(ev.Kv.Value) {
					broker.leader = string(ev.Kv.Value)
					broker.votechan <- struct{}{}
				}
			case "DELETE":
				go broker.vote()
			}
		}
	}
}

func (broker *Broker) vote() {
	logs.Info("I am voting")
	select {
	case <-broker.votechan:
		return
	case <-time.After(time.Duration(rand.New(rand.NewSource(time.Now().Unix())).Intn(200)) * time.Millisecond):
		resp, err := broker.Etcd.Client.Grant(context.TODO(), 10)
		if err != nil {
			logs.Error(err)
			os.Exit(1)
		}
		if txnresp, _ := broker.Etcd.Client.Txn(context.TODO()).
			If(clientv3.Compare(clientv3.CreateRevision("leader"), "=", 0)).
			Then(clientv3.OpPut("leader", "broker"+strconv.Itoa(broker.id), clientv3.WithLease(resp.ID))).
			Commit(); txnresp.Succeeded {
			go broker.leaderhearbeat(resp)
			go broker.watchmembers()
		}
	}
}

func (broker *Broker) leaderhearbeat(resp *clientv3.LeaseGrantResponse) {
	for {
		<-time.After(time.Second * 8)
		logs.Info("leaderhearbeat")
		_, err := broker.Etcd.Client.KeepAliveOnce(context.TODO(), resp.ID)
		if err != nil {
			logs.Error(err)
			os.Exit(1)
		}
	}
}

func (broker *Broker) getmembers() {
	var tmp []string
	broker.members = tmp
	resp, _ := broker.Etcd.Client.Get(context.TODO(), "broker", clientv3.WithPrefix())
	for _, v := range resp.Kvs {
		broker.members = append(broker.members, string(v.Value))
	}
	logs.Info("all brokers:")
	for k, v := range broker.members {
		logs.Info(k, v)
	}
}

func (broker *Broker) watchmembers() {
	broker.getmembers()
	wch := broker.Etcd.Client.Watch(context.TODO(), "broker", clientv3.WithPrefix())
	for wresp := range wch {
		for _, ev := range wresp.Events {
			switch ev.Type.String() {
			case "PUT":
				logs.Info("creat broker:", string(ev.Kv.Value))
				broker.members = append(broker.members, string(ev.Kv.Value))
				logs.Info("all brokers:")
				for k, v := range broker.members {
					logs.Info(k, v)
				}
			case "DELETE":
				broker.getmembers()
			}
		}
	}
}
