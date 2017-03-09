package msg

import (
	"strconv"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/golang/protobuf/proto"
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

func (topic *topic) load() {
	ro := gorocksdb.NewDefaultReadOptions()
	defer ro.Destroy()
	ro.SetFillCache(false)
	it := topic.rocks.NewIterator(ro)
	defer it.Close()
	it.SeekToFirst()
	for it = it; it.Valid(); it.Next() {
		tmp := &Msg{}
		err := proto.Unmarshal(it.Value().Data(), tmp)
		if err != nil {
			logs.Error(err)
		}
		topic.put(tmp.GetAlivetime(), tmp.GetBody())
		it.Key().Free()
		it.Value().Free()
	}
	if err := it.Err(); err != nil {
		logs.Error(err)
	}
}

func (topic *topic) put(alivetime, body string, offset ...string) {
	if len(offset) == 0 {
		tmp := newMessage(alivetime, body)
		topic.mutex.Lock()
		topic.msg[offset[0]] = tmp
		topic.mutex.Unlock()
		topic.dump(offset[0], tmp)
	} else {
		tmp := newMessage(alivetime, body)
		topic.mutex.Lock()
		topic.msg[strconv.Itoa(len(topic.msg))] = tmp
		topic.mutex.Unlock()
		topic.dump(strconv.Itoa(len(topic.msg)), tmp)
	}
}

func (topic *topic) get(offset string) *message {
	topic.mutex.RLock()
	message := topic.msg[offset]
	topic.mutex.Unlock()
	return message
}

func (topic *topic) dump(offset string, msg *message) {
	body := &Msg{
		Alivetime: proto.String(msg.alivetime),
		Body:      proto.String(msg.body),
	}
	data, err := proto.Marshal(body)
	if err != nil {
		logs.Error(err)
	}
	wo := gorocksdb.NewDefaultWriteOptions()
	defer wo.Destroy()
	err = topic.rocks.Put(wo, []byte(offset), []byte(data))
	if err != nil {
		logs.Error(err)
	}
}

func (topic *topic) del(offset string) {
	topic.mutex.Lock()
	delete(topic.msg, offset)
	topic.mutex.Unlock()
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
