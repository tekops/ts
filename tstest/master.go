package tstest

import (
	"fmt"
	"github.com/tekops/ts/master"
	//"github.com/tekops/ts/utils"
	//"log"
	"net"
)

func newLocalListener(addr string) net.Listener {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(fmt.Sprintf("tstest: failed to listen on a port: %v", err))
	}
	return l
}

// A Master defines parameters for running a Master
type Master struct {
	//CCListner     net.Listener
	//MinionListner net.Listener
	RMaster *master.Master
}

//NewUnstartedMaster returns a new master but does not starts it
func NewUnstartedMaster() *Master {
	config := master.ReadConfig("/home/manoj/gocode/src/github.com/tekops/ts/utils/conf/master.yaml")
	fmt.Printf("Config: %#v \n", config)

	//miaddr := "127.0.0.1:2002"
	//claddr := "127.0.0.1:2001"
	return &Master{
		RMaster: master.New(config),
		//CCListner:     newLocalListener(claddr),
		//MinionListner: newLocalListener(miaddr),
	}
}

// NewServer starts and returns a new Server.
// The caller should call Close when finished, to shut it down.
func NewMaster() *Master {
	nm := NewUnstartedMaster()
	nm.Start()
	return nm
}

func (master *Master) Start() {
	//go master.RMaster.ServeMinions(master.MinionListner)
	//go master.RMaster.ServeCC(master.CCListner)
}

func (master *Master) Close() {
	//master.CCListner.Close()
	//master.MinionListner.Close()
}
