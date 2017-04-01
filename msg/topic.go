package msg

import (
	"strconv"
	"sync"
	"time"

	"errors"

	"github.com/astaxie/beego/logs"
	"github.com/golang/protobuf/proto"
	"github.com/ohmq/ohmyqueue/inrpc"
	"github.com/tecbot/gorocksdb"
)

type topic struct {
	rocks  *gorocksdb.DB
	mutex  *sync.RWMutex
	offset int64
	msg    map[int64]*message
}

func newTopic(name string) *topic {
	rocks := newDB(name)
	msg := make(map[int64]*message)
	return &topic{
		rocks:  rocks,
		mutex:  new(sync.RWMutex),
		offset: 0,
		msg:    msg,
	}
}

func (topic *topic) load() {
	logs.Info("load form local")
	ro := gorocksdb.NewDefaultReadOptions()
	defer ro.Destroy()
	ro.SetFillCache(false)
	it := topic.rocks.NewIterator(ro)
	defer it.Close()
	it.SeekToFirst()
	for ; it.Valid(); it.Next() {
		tmp := &Msg{}
		err := proto.Unmarshal(it.Value().Data(), tmp)
		if err != nil {
			logs.Error(err)
		}
		topic.mutex.Lock()
		off, _ := strconv.ParseInt(string(it.Key().Data()), 10, 64)
		topic.msg[off] = newMessage(tmp.GetAlivetime(), tmp.GetBody())
		topic.mutex.Unlock()
		if off > topic.offset {
			topic.offset = off
		}
		//it.Key().Free()
		//it.Value().Free()
	}
	if err := it.Err(); err != nil {
		logs.Error(err)
	}
}

func (topic *topic) put(alivetime, body string, offset ...int64) (off int64) {
	if len(offset) != 0 {
		tmp := newMessage(alivetime, body)
		topic.mutex.Lock()
		topic.msg[offset[0]] = tmp
		topic.mutex.Unlock()
		topic.dump(offset[0], tmp)
		off = offset[0]
	} else {
		tmp := newMessage(alivetime, body)
		topic.mutex.Lock()
		topic.msg[topic.offset] = tmp
		topic.offset++
		topic.mutex.Unlock()
		topic.dump(int64(len(topic.msg)), tmp)
		off = topic.offset - 1
	}
	return
}

func (topic *topic) get(offset int64) (off int64, message string, err error) {
	topic.mutex.RLock()
	defer topic.mutex.RUnlock()
	for ; offset <= topic.offset-1; offset++ {
		if v, ok := topic.msg[offset]; ok {
			off = offset
			message = v.body
			err = nil
			return
		}
	}
	off = -1
	err = errors.New("not exits")
	return
}

func (topic *topic) dump(offset int64, msg *message) {
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
	err = topic.rocks.Put(wo, []byte(strconv.FormatInt(offset, 10)), []byte(data))
	if err != nil {
		logs.Error(err)
	}
}

func (topic *topic) del(offset int64) {
	topic.mutex.Lock()
	delete(topic.msg, offset)
	topic.mutex.Unlock()
	wo := gorocksdb.NewDefaultWriteOptions()
	defer wo.Destroy()
	err := topic.rocks.Delete(wo, []byte(strconv.FormatInt(offset, 10)))
	if err != nil {
		logs.Error(err)
	}
}

func (topic *topic) clean() {
	for {
		<-time.After(time.Minute)
		for k, v := range topic.msg {
			logs.Info("clean start", k, v)
			if v.alivetime <= strconv.FormatInt(time.Now().Unix(), 10) {
				logs.Info("del", k, v)
				topic.del(k)
			}
		}
	}
}

func (topic *topic) getall() []*inrpc.Msg {
	var msgs = make([]*inrpc.Msg, 100)
	topicname := topic.rocks.Name()[10:13]
	for offset, body := range topic.msg {
		msgs = append(msgs, &inrpc.Msg{
			Topic:     topicname,
			Offset:    offset,
			Alivetime: body.alivetime,
			Body:      body.body,
		})
	}
	return msgs
}
