package server

import (
	"context"
	"net"

	"strconv"

	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/ohmq/ohmyqueue/etcd"
	"github.com/ohmq/ohmyqueue/msg"
)

type Broker struct {
	id   int
	etcd *etcd.Etcd
	ip   string
	port int
	msgs *msg.Msgs
}

func NewBroker(id int, msgs *msg.Msgs) *Broker {
	ip := ""
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	return &Broker{
		id:   id,
		etcd: etcd.NewEtcd(),
		ip:   ip,
		port: 9988,
		msgs: msgs,
	}
}

func (broker *Broker) Start() {
	go broker.etcd.Heartbeat("broker/index"+strconv.Itoa(broker.id), broker.ip+":"+strconv.Itoa(broker.port), 10)
	broker.watch()
}

func (broker *Broker) watch() {
	wch := broker.etcd.Client.Watch(context.TODO(), "broker/index"+strconv.Itoa(int(broker.id))+"/topics", clientv3.WithPrefix())
	for wresp := range wch {
		for _, ev := range wresp.Events {
			switch ev.Type.String() {
			case "PUT":
				broker.msgs.Update(strings.Split(string(ev.Kv.Value), ","))
			}
		}
	}
}
