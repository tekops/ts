package cc

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/tekops/ts/utils"
	"io/ioutil"
	"launchpad.net/goyaml"
	"net"
)

var ()

type Config struct {
	Masters []string
}

type CC struct {
	Masters []string
}

func New(config Config) *CC {
	return &CC{
		Masters: config.Masters,
	}
}

func (c CC) SendCommnd(target, dest, command string) {
	conn, err := net.Dial("tcp", dest)
	if err != nil {
		glog.Errorln(err)
		return
	}
	defer conn.Close()
	InCh := make(chan *utils.Message, 1024)

	msg := utils.NewMessage()
	msg.Cmd = command
	msg.JobType = utils.CMD_JOBTYPE
	msg.Target = target
	log_prefix := fmt.Sprintf("cc to %v", dest)
	utils.SendMessage(log_prefix, conn, &msg)
	go utils.ReadFromConn("cc ", conn, InCh)
	msg1 := <-InCh
	glog.Infoln("got from master", msg1)
	conn.Close()
	//    for msg := range InCh {
	//		glog.Infoln("got from master", msg)
	//	}

}

func ReadConfig(cfgFile string) Config {
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
