package msg

import (
	"github.com/astaxie/beego/logs"
	"github.com/ohmq/ohmyqueue/config"
	"github.com/tecbot/gorocksdb"
)

func newDB(name string) *gorocksdb.DB {
	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(gorocksdb.NewLRUCache(3 << 30))
	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(opts, config.Conf.Omq.Logdir+name+".db")
	logs.Info(config.Conf.Omq.Logdir + name + ".db")
	if err != nil {
		logs.Info(err)
	}
	return db
}
