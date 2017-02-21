package msg

import (
	"fmt"
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
			msgs.AddTopic(name, Topic{Alivetime: int64(time.Second * 120), Message: make(map[string]Msg)})
		}
	}
	for name, topic := range msgs.Topics {
		fmt.Printf("%s %#v \n", name, topic)
	}
}

func (msgs *Msgs) Put(topic string, msg Msg) {
	log.Println("put")
	log.Println(msgs.Topics[topic])
	msgs.Topics[topic].Message[strconv.Itoa(len(msgs.Topics[topic].Message))] = msg
	for k, v := range msgs.Topics[topic].Message {
		log.Printf("%s %#v\n", k, v)
	}
	log.Println(msgs.Topics[topic])

}

func (msgs *Msgs) Get(topic string, offset string) Msg {
	return msgs.Topics[topic].Message[offset]
}
