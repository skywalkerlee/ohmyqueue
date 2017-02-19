package etcd

import (
	"context"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
)

type Etcd struct {
	Client *clientv3.Client
}

func NewEtcd() *Etcd {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"172.27.31.156:12379", "172.27.31.156:22379", "172.27.31.156:32379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return &Etcd{cli}
}

func (etcd *Etcd) Heartbeat(key, value string, timeout int64) {
	resp, err := etcd.Client.Grant(context.TODO(), timeout)
	if err != nil {
		log.Fatal(err)
	}
	_, err = etcd.Client.Put(context.TODO(), key, value, clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-time.After(time.Second):
			_, err = etcd.Client.KeepAliveOnce(context.TODO(), resp.ID)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
