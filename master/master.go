package master

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/tekops/ts/utils"
	"io/ioutil"
	"launchpad.net/goyaml"
	"net"
	// "time"
)

var (
	InCh   = make(chan *utils.Message, 1024)
	CCInCh = make(chan *utils.Message, 1024)
)

// A Master defines parameters for running a Master
type Master struct {
	MLAddr        string //TCP address for communication with minions
	CLAddr        string //TCP address for Command Relay
	ActiveMinions *utils.NodeList
	ActiveCC      *utils.NodeList
}

type MasterConfig struct {
	ListenAddr        string "listen_addr"
	NodeId            string "node_id"
	NodeUUID          string "node_uuid"
	CommandCenterAddr string "command_center_addr"
}

// Instantiate Master
func New(conf MasterConfig) *Master {
	return &Master{
		MLAddr:        conf.ListenAddr,
		CLAddr:        conf.CommandCenterAddr,
		ActiveMinions: utils.NewNodeList(),
		ActiveCC:      utils.NewNodeList(),
	}
}

func (m *Master) Start() {
	ml, err := net.Listen("tcp", m.MLAddr)
	if err != nil {
		glog.Fatalln(err)
	}
	go m.ServeMinions(ml)
	cl, err := net.Listen("tcp", m.CLAddr)
	if err != nil {
		glog.Fatalln(err)
	}
	go m.ServeCC(cl)
	go m.handleCC()
	go m.handleMinions()
	glog.Infoln("Master Started")
}

//listen for minions
func (m *Master) ServeMinions(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			glog.Errorln("Error in accepting connections from minions: ", err)
			continue
		}
		remoteAddr := conn.RemoteAddr().String()
		m.ActiveMinions.Add(remoteAddr, conn)
		go utils.ReadFromConn("master", conn, InCh)

		if glog.V(3) {
			glog.Infoln("Starting to accept data from minion", remoteAddr)
			glog.Infof("ActiveMinions: %#v \n", m.ActiveMinions.List())
		}
	}
}

//listen for minions an
func (m *Master) ServeCC(l net.Listener) {

	for {
		conn, err := l.Accept()
		if err != nil {
			glog.Errorln("Error in accepting connections from minions: ", err)
			continue
		}
		remoteAddr := conn.RemoteAddr().String()
		m.ActiveCC.Add(remoteAddr, conn)
		go utils.ReadFromConn("master", conn, CCInCh)
		if glog.V(3) {
			glog.Infoln("Starting to accept Commands from ", remoteAddr)
			glog.Infof("ActiveMinions: %#v \n", m.ActiveMinions.List())
		}
	}
}

//read from chan and echo massage
func (m *Master) handleMinions() {
	for msg := range InCh {
		if glog.V(3) {
			glog.Infof("Got from minion %#v\n", msg)
		}
		switch msg.JobType {
		case utils.CMD_JOBTYPE:
			if glog.V(2) {
				glog.Infoln("got command answer from minion", msg.JobId)
				glog.Infof("activeCC: %#v", m.ActiveCC.List())
			}

			for _, cc := range m.ActiveCC.List() {
				if cc.IPADDR == msg.IPADDR {
					utils.SendMessage("master", cc.Conn, msg)
				}
			}
		}
	}
}

//read from chan and echo massage
func (m *Master) handleCC() {
	for msg := range CCInCh {
		msg.JobId = utils.UUID()
		if glog.V(2) {
			glog.Infof("Got from CC %#v\n", msg.JobId)
		}
		if glog.V(3) {
			glog.Infof("Got from CC %#v\n", msg)
		}
		target := ParseTarget(msg.Target)
		switch target {
		case "*":
			for _, minion := range m.ActiveMinions.List() {
				log_prefix := fmt.Sprintf("master to %v", minion.IPADDR)
				utils.SendMessage(log_prefix, minion.Conn, msg)
				if glog.V(2) {
					glog.Infoln("Message sent to minion", msg.JobId)
				}
			}
		default:
			glog.Infoln("Unknown Case")
		}
	}
}

func ReadConfig(cfgFile string) MasterConfig {
	var config MasterConfig
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

func ParseTarget(targets string) string {
	return "*"
}
