package broker

import (
	"os"

	"github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/config"
	"github.com/tecbot/gorocksdb"
)

func newDB() *gorocksdb.DB {
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))
	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(opts, config.Conf.Omq.Logdir+"omq.db")
	logs.Info(config.Conf.Omq.Logdir + "omq.db")
	if err != nil {
		logs.Info(err)
		os.Exit(1)
	}
	return db
}

// func (broker *Broker) load() {
// 	ro := gorocksdb.NewDefaultReadOptions()
// 	defer ro.Destroy()
// 	ro.SetFillCache(false)
// 	it := broker.db.NewIterator(ro)
// 	defer it.Close()
// 	it.SeekToFirst()
// 	for it = it; it.Valid(); it.Next() {
// 		broker.msg[string(it.Key().Data())] = string(it.Value().Data())
// 		it.Key().Free()
// 		it.Value().Free()
// 	}
// 	if err := it.Err(); err != nil {
// 		logs.Error(err)
// 	}
// }

// func (broker *Broker) dump(key, value string) {
// 	wo := gorocksdb.NewDefaultWriteOptions()
// 	defer wo.Destroy()
// 	err := broker.db.Put(wo, []byte(key), []byte(value))
// 	if err != nil {
// 		logs.Error(err)
// 	}
// }

// func (broker *Broker) clean(key string) {
// 	wo := gorocksdb.NewDefaultWriteOptions()
// 	defer wo.Destroy()
// 	err := broker.db.Delete(wo, []byte(key))
// 	if err != nil {
// 		logs.Error(err)
// 	}
// }
