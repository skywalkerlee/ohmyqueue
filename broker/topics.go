package broker

import (
	"context"
	"math/rand"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/clientv3"
)

func (broker *Broker) getTopics() {
	resp, _ := broker.Client.Get(context.TODO(), "topic", clientv3.WithPrefix())
	for _, v := range resp.Kvs {
		broker.topics = append(broker.topics, string(v.Key))
	}
}

func (broker *Broker) watchTopics() {
	resp, _ := broker.Client.Get(context.TODO(), "topic", clientv3.WithPrefix())
	for _, v := range resp.Kvs {
		broker.topics = append(broker.topics, string(v.Key))
	}
	wch := broker.Client.Watch(context.TODO(), "topic", clientv3.WithPrefix())
	for wresp := range wch {
		for _, ev := range wresp.Events {
			switch ev.Type.String() {
			case "PUT":
				logs.Info("creat topic:", string(ev.Kv.Value))
				broker.topics = append(broker.topics, string(ev.Kv.Key))
			}
		}
	}
}

func (broker *Broker) watchTopicLeader(name string) {
	resp, _ := broker.Client.Get(context.TODO(), "topic"+name+"leader")
	if resp.Count == 0 {
		go broker.vote()
	} else {
		logs.Info("topic"+name+"leader is:", string(resp.Kvs[0].Value))
	}
	wch := broker.Client.Watch(context.TODO(), "topic"+name+"leader")
	for wresp := range wch {
		for _, ev := range wresp.Events {
			switch ev.Type.String() {
			case "PUT":
				if string(ev.Kv.Value) == broker.ip+":"+broker.clientport {
					broker.leaders = append(broker.leaders, name)
				}
			}
		}
	}
}

func (broker *Broker) voteTopicleader(name string) {
	logs.Info("I am voting")
	<-time.After(time.Duration(rand.New(rand.NewSource(time.Now().Unix())).Intn(200)) * time.Millisecond)
	resp, err := broker.Client.Grant(context.TODO(), 5)
	if err != nil {
		logs.Error(err)
	}
	if txnresp, _ := broker.Client.Txn(context.TODO()).
		If(clientv3.Compare(clientv3.CreateRevision("topic"+name+"leader"), "=", 0)).
		Then(clientv3.OpPut("topic"+name+"leader", broker.ip+":"+broker.clientport, clientv3.WithLease(resp.ID))).
		Commit(); txnresp.Succeeded {
		go broker.leaderhearbeat(resp)
	}
}

func (broker *Broker) topicLeaderHeartbeat(resp *clientv3.LeaseGrantResponse) {
	for {
		<-time.After(time.Second * 4)
		logs.Info("leaderhearbeat")
		_, err := broker.Client.KeepAliveOnce(context.TODO(), resp.ID)
		if err != nil {
			logs.Error(err)
		}
	}
}
