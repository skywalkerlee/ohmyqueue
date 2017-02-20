package msg

import (
	"strconv"
)

type Topic struct {
	Message map[string]Msg
}

type Msgs struct {
	Topics map[string]Topic
}

func NewMsgs() *Msgs {
	return &Msgs{}
}

func (msgs *Msgs) AddTopic(name string, topic Topic) {
	msgs.Topics[name] = topic
}

func (msgs *Msgs) DelTopic(name string) {
	delete(msgs.Topics, name)
}

func (msgs *Msgs) Update(names []string) {
	for _, name := range names {
		if _, ok := msgs.Topics[name]; !ok {
			msgs.Topics[name] = Topic{}
		}
	}
}

func (msgs *Msgs) Put(topic string, msg Msg) {
	msgs.Topics[topic].Message[strconv.Itoa(len(msgs.Topics[topic].Message))] = msg
}

func (msgs *Msgs) Get(topic string, offset string) Msg {
	return msgs.Topics[topic].Message[offset]
}
