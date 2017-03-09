package msg

import (
	"strconv"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/tecbot/gorocksdb"
)

type topic struct {
	rocks *gorocksdb.DB
	mutex *sync.RWMutex
	msg   map[string]*message
}

func newTopic(name string) *topic {
	rocks := newDB(name)
	msg := make(map[string]*message)
	return &topic{
		rocks: rocks,
		mutex: new(sync.RWMutex),
		msg:   msg,
	}
}

func (topic *topic) put(alivetime, body string) {
	topic.mutex.Lock()
	topic.msg[strconv.Itoa(len(topic.msg))] = newMessage(alivetime, body)
	topic.mutex.Unlock()
	topic.dump(strconv.Itoa(len(topic.msg)), body)
}

func (topic *topic) get(offset string) *message {
	topic.mutex.RLock()
	message := topic.msg[offset]
	topic.mutex.Unlock()
	return message
}

func (topic *topic) dump(offset, body string) {
	wo := gorocksdb.NewDefaultWriteOptions()
	defer wo.Destroy()
	err := topic.rocks.Put(wo, []byte(offset), []byte(body))
	if err != nil {
		logs.Error(err)
	}
}

func (topic *topic) del(offset string) {
	topic.mutex.Lock()
	delete(topic.msg, offset)
	topic.mutex.Lock()
	wo := gorocksdb.NewDefaultWriteOptions()
	defer wo.Destroy()
	err := topic.rocks.Delete(wo, []byte(offset))
	if err != nil {
		logs.Error(err)
	}
}

func (topic *topic) clean() {
	for {
		<-time.After(time.Minute)
		for k, v := range topic.msg {
			if v.alivetime >= strconv.FormatInt(time.Now().Unix(), 10) {
				topic.del(k)
			}
		}
	}
}
