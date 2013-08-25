package utils

import (
	"net"
	"time"
)

type Message struct {
	JobId     string
	JobType   int
	Cmd       string
	OutPut    string
	NodeId    string
	Conn      net.Conn
	TimeStamp time.Time
	Target    string
	IPADDR    string
	//NodeUUID  string
}

func NewMessage() Message {
	msg := Message{
		TimeStamp: time.Now(),
	}
	return msg
}

const (
	CMD_JOBTYPE = iota
	MON_JOB_TYPE
)
