package utils

import (
	"github.com/golang/glog"
	"io/ioutil"
	"launchpad.net/goyaml"
)

type Config interface{}

func ReadConfig(cfgFile string) interface{} {
	var config Config
	if cfgFile == "" {
		cfgFile = "conf/config.yaml"
	}
	file, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		glog.Fatal("error reading config file ", err)
	}
	err = goyaml.Unmarshal(file, &config)
	if err != nil {
		glog.Fatal("error decoding config file ", err)
	}
	return config
}
