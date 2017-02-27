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
		Timeout    int    `toml:"timeout"`
		Logdir     string `toml:"logdir"`
	} `toml:"omq"`
}

var Conf Config

func init() {
	if _, err := toml.DecodeFile("./omq.conf", &Conf); err != nil {
		logs.Error(err)
		os.Exit(1)
	}
}
