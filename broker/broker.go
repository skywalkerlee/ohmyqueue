package broker

import (
	"net"

	"github.com/coreos/etcd/clientv3"
	"github.com/ohmq/ohmyqueue/etcd"
)

type Broker struct {
	id         int
	Client     *clientv3.Client
	ip         string
	clientport string
	innerport  string
	votechan   chan struct{}
	leader     string
	members    map[string]string
	topics     []string
	leaders    []string
	// ldtopic    map[string]map[string]string
	// mbtopic    map[string]map[string]string
}

func NewBroker(id int, cliport string, inport string) *Broker {
	var ip string
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	var ts []string
	var ls []string
	return &Broker{
		id:         id,
		Client:     etcd.NewEtcd().Client,
		ip:         ip,
		members:    make(map[string]string),
		clientport: cliport,
		innerport:  inport,
		topics:     ts,
		leaders:    ls,
	}
}

func (broker *Broker) Close() {
	broker.Client.Close()
}
