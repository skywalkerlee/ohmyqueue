package msg

import (
	"github.com/ohmq/ohmyqueue/inrpc"
)

type Topics map[string]*topic

func NewTopics() Topics {
	topics := make(map[string]*topic)
	return topics
}

func (topics Topics) AddTopic(name string) {
	topics[name] = newTopic(name)
	topics[name].load()
	go topics[name].clean()
}

func (topics Topics) Put(topic, alivetime, body string, offset ...string) {
	topics[topic].put(alivetime, body, offset...)
}

func (topics Topics) Get(topic, offset string) string {
	return topics[topic].get(offset).body
}

func (topics Topics) GetAll(topic string) []*inrpc.Msg {
	return topics[topic].getall()
}
