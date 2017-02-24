package etcd

import (
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
