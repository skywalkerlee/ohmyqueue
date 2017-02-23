package broker

import (
	"net"

	"errors"

	"github.com/ohmq/ohmyqueue/etcd"
)

type Broker struct {
	id         int
	Etcd       *etcd.Etcd
	ip         string
	clientport int
	innerport  int
	votechan   chan struct{}
	// Msgs       *msg.Msgs
	leader  string
	members []string
	msg     map[string]string
}

func NewBroker(id int, cliport int, inport int) *Broker {
	ip := ""
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}
	var members []string
	return &Broker{
		id:         id,
		Etcd:       etcd.NewEtcd(),
		ip:         ip,
		members:    members,
		clientport: cliport,
		innerport:  inport,
		msg:        make(map[string]string),
	}
}

func (broker *Broker) Put(offset, body string) error {
	if _, ok := broker.msg[offset]; ok {
		broker.msg[offset] = body
		return nil
	}
	return errors.New("offset is allready exist now")
}

func (broker *Broker) Get(offset string) (string, error) {
	if v, ok := broker.msg[offset]; ok {
		return v, nil
	}
	return "", errors.New("offset is not exist")
}
