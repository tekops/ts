package tstest

import (
	"flag"
	"github.com/tekops/ts/cc"
	"github.com/tekops/ts/master"
	"github.com/tekops/ts/minion"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	config := master.ReadConfig("/home/manoj/gocode/src/github.com/tekops/ts/conf/master.yaml")
	master := master.New(config)
	go master.Start()
	time.Sleep(1 * time.Second)
	minionConfig := minion.ReadConfig("/home/manoj/gocode/src/github.com/tekops/ts/conf/minion.yaml")
	minion := minion.New(minionConfig)
	go minion.Start()
	time.Sleep(1 * time.Second)
	ccConfig := cc.ReadConfig("/home/manoj/gocode/src/github.com/tekops/ts/conf/cc.yaml")
	cc := cc.New(ccConfig)
	go cc.SendCommnd("*", master.CLAddr, "df -h")
	go cc.SendCommnd("*", master.CLAddr, "uname -a")
	time.Sleep(60 * time.Second)
}

func init() {
	flag.Parse()
	flag.Set("logtostderr", "true")
	flag.Set("v", "2")
}
