package msg

import (
	"log"
	"strconv"
	"time"
)

type Topic struct {
	Alivetime int64
	Message   map[string]Msg
}

type Msgs struct {
	Topics map[string]Topic
}

func NewMsgs() *Msgs {
	return &Msgs{Topics: make(map[string]Topic)}
}

func (msgs *Msgs) AddTopic(name string, topic Topic) {
	msgs.Topics[name] = topic
}

func (msgs *Msgs) DelTopic(name string) {
	delete(msgs.Topics, name)
}

func (msgs *Msgs) Update(names []string) {
	log.Println("Update")
	for _, name := range names {
		if _, ok := msgs.Topics[name]; !ok {
			msgs.Topics[name] = Topic{Alivetime: int64(time.Second * 120), Message: make(map[string]Msg)}
		}
	}
}

func (msgs *Msgs) Put(topic string, msg Msg) {
	log.Println("put")
	index := strconv.Itoa(len(msgs.Topics[topic].Message))
	log.Println(index)
	aa := msgs.Topics[topic]
	aa.Message[index] = msg
	for k, v := range msgs.Topics[topic].Message {
		log.Println(k, v)
	}
}

func (msgs *Msgs) Get(topic string, offset string) Msg {
	return msgs.Topics[topic].Message[offset]
}
