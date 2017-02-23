package etcd

import (
	"context"
	"time"

	"os"

	log "github.com/astaxie/beego/logs"
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
		log.Error(err)
		os.Exit(1)
	}
	return &Etcd{cli}
}

func (etcd *Etcd) Heartbeat(key, value string, timeout int64) {
	resp, err := etcd.Client.Grant(context.TODO(), timeout)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	_, err = etcd.Client.Put(context.TODO(), key, value, clientv3.WithLease(resp.ID))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	for {
		select {
		case <-time.After(time.Second * 8):
			log.Info("hearbeat")
			_, err = etcd.Client.KeepAliveOnce(context.TODO(), resp.ID)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
		}
	}
}
