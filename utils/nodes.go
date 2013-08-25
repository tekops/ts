package utils

import (
	"net"
	"sync"
	"time"
)

//Node is represents physical servers(nodes)
type Node struct {
	IPADDR   string
	Conn     net.Conn
	NodeId   string
	LastSeen time.Time
	NodeUUID string
}

func NewNode(addr string, c net.Conn) (Node, bool) {
	node := Node{
		IPADDR:   addr,
		Conn:     c,
		LastSeen: time.Now(),
	}
	return node, true
}

type NodeList struct {
	nodes map[string]Node
	mu    sync.RWMutex
}

func NewNodeList() *NodeList {
	nl := &NodeList{nodes: make(map[string]Node)}
	return nl
}

//Add Node to list if its node is present change conn to newer one
func (n *NodeList) Add(addr string, conn net.Conn) (Node, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	//ipAddr, _, _ := net.SplitHostPort(addr)
	ipAddr := addr
	if node, ok := n.nodes[ipAddr]; ok {
		node.Conn = conn
		return node, true
	}
	node, ok := NewNode(ipAddr, conn)
	if ok {
		n.nodes[ipAddr] = node
		return node, true
	}
	return node, false
}

func (n *NodeList) List() map[string]Node {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.nodes
}
