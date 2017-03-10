package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/astaxie/beego/logs"
)

type Config struct {
	Etcd struct {
		Addr []string `toml:"addr"`
	} `toml:"etcd"`
	Omq struct {
		Index      int    `toml:"index"`
		Clientport string `toml:"clientport"`
		Innerport  string `toml:"innerport"`
		Timeout    int64  `toml:"timeout"`
		Logdir     string `toml:"logdir"`
	} `toml:"omq"`
	Topic struct {
		Alivetime int64 `toml:"alivetime"`
	} `toml:"topic"`
}

var Conf Config

func init() {
	logs.SetLogger("console")
	logs.SetLogger(logs.AdapterFile, `{"filename":"`+Conf.Omq.Logdir+`omq.log"}`)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	if _, err := toml.DecodeFile("./omq.conf", &Conf); err != nil {
		logs.Error(err)
		os.Exit(1)
	}
}
