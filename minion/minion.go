package minion

import (
	"github.com/golang/glog"
	"github.com/tekops/ts/utils"
	"io/ioutil"
	"launchpad.net/goyaml"
	"net"
	"os/exec"
	"strings"
)

var (
	inch  = make(chan *utils.Message, 1024)
	outch = make(chan *utils.Message, 1024)
)

//active masters
var (
	activeMasters = utils.NewNodeList()
)

type MinionConfig struct {
	Masters []string
}

type Minion struct {
	ActiveMasters *utils.NodeList
	Masters       []string
	ID            string
}

func New(conf MinionConfig) *Minion {
	return &Minion{
		Masters: conf.Masters,
	}
}

func (m *Minion) Start() {
	for _, master := range m.Masters {
		conn, err := net.Dial("tcp", master)
		if err != nil {
			glog.Errorln(err)
			continue
		}
		if glog.V(2) {
			glog.Infoln("connected to master: ", master)
		}
		go utils.ReadFromConn("minion", conn, inch)
		go m.handleMaster()
	}
}

//read from chan and echo massage
func (m *Minion) handleMaster() {
	for msg := range inch {
		if glog.V(2) {
			glog.Infoln("Got from master", msg.JobId)
		}
		switch msg.JobType {
		case utils.CMD_JOBTYPE:
			go m.executeCommand(msg)

		}
	}
}

func ReadConfig(cfgFile string) MinionConfig {
	var config MinionConfig
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
	if len(config.Masters) == 0 {
		glog.Errorln("No Masters?")
	}
	return config
}

func (minion *Minion) executeCommand(msg *utils.Message) {
	str := strings.Fields(msg.Cmd)
	cmdName := str[0]
	aname, err := exec.LookPath(cmdName)
	if err != nil {
		aname = cmdName
	}

	args := str[1:]
	cmd := exec.Cmd{
		Path: aname,
		Args: append([]string{cmdName}, args...),
	}
	if glog.V(10) {
		glog.Infof("Command is: %#v\n", cmd)
	}
	out, err := cmd.Output()
	if err != nil {
		glog.Errorln("error executing ", err)

	} else {
		if glog.V(10) {
			glog.Infoln(string(out))
		}
		returnMsg := &utils.Message{
			JobId:  msg.JobId,
			Cmd:    msg.Cmd,
			OutPut: string(out),
			NodeId: minion.ID,
			Conn:   msg.Conn,
		}

		utils.SendMessage("minion to master", msg.Conn, returnMsg)
	}
}
