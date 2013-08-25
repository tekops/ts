package utils

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"net"
	"os"
)

//reads from connection encodes & sends to chan
func ReadFromConn(caller string, conn net.Conn, ch chan *Message) {
	if glog.V(3) {
		glog.Infoln(caller, ":got connection from:", conn.RemoteAddr().String())
	}
	dec := json.NewDecoder(conn)
	for {
		var m *Message
		err := dec.Decode(&m)
		if err != nil {
			if glog.V(3) {
				glog.Errorln(caller, ":", err)
			}
			break
		}
		if glog.V(3) {
			glog.Infof("%v: got %#v \n", caller, m)
		}
		m.Conn = conn
		ch <- m
	}
}

func SendMessage(caller string, conn net.Conn, msg *Message) {
	msg.Conn = nil
	enc := json.NewEncoder(conn)
	err := enc.Encode(msg)
	if err != nil {
		glog.Errorln("Error encoding data ", err)
	}
	glog.Infoln(caller, ": Data Send", msg.JobId)
}

func UUID() string {
	f, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
