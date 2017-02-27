package etcd

import (
	"time"

	"os"

	log "github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/clientv3"
	"github.com/ohmq/ohmyqueue/config"
)

type Etcd struct {
	Client *clientv3.Client
}

func NewEtcd() *Etcd {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Conf.Etcd.Addr,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	return &Etcd{cli}
}
