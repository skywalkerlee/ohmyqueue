package broker

import (
	"net"

	"errors"

	"sync"

	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/clientv3"
	"github.com/ohmq/ohmyqueue/etcd"
	"github.com/ohmq/ohmyqueue/inrpc"
	"github.com/tecbot/gorocksdb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	msg        map[string]string
	lock       *sync.Mutex
	db         *gorocksdb.DB
	topics     []string
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
	return &Broker{
		id:         id,
		Client:     etcd.NewEtcd().Client,
		ip:         ip,
		members:    make(map[string]string),
		clientport: cliport,
		innerport:  inport,
		msg:        make(map[string]string),
		lock:       new(sync.Mutex),
		db:         newDB(),
		topics:     ts,
	}
}

func (broker *Broker) Close() {
	broker.Client.Close()
	broker.db.Close()
}

func (broker *Broker) Put(body string) {
	logs.Info("put")
	broker.msg[strconv.Itoa(len(broker.msg))] = body
	broker.dump(strconv.Itoa(len(broker.msg)), body)
	if len(broker.members) == 0 {
		return
	}
	for _, v := range broker.members {
		go func(addr string) {
			defer func() {
				if err := recover(); err != nil {
					logs.Error(err)
				}
			}()
			broker.sync(addr)
		}(v)
	}
}

func (broker *Broker) Del(offset string) {
	delete(broker.msg, offset)
}

func (broker *Broker) Len() int {
	return len(broker.msg)
}

func (broker *Broker) Get(offset string) (string, error) {
	if v, ok := broker.msg[offset]; ok {
		return v, nil
	}
	return "", errors.New("offset is not exist")
}

func (broker *Broker) sync(addr string) {
	logs.Info("sync", addr)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err.Error())
	}
	client := inrpc.NewInClient(conn)
	stream, err := client.SyncMsg(context.TODO())
	if err != nil {
		panic(err.Error())
	}
	for k, v := range broker.msg {
		stream.Send(&inrpc.Msg{Topic: "test", Offset: k, Body: v})
	}
	statuscode, err := stream.CloseAndRecv()
	if err != nil {
		panic(err.Error())
	}
	logs.Info(statuscode.GetSum())
}
