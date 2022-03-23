package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	Remove = 1
	Add    = 2
)

type Node struct {
	id   int
	addr string
}

func (s *Node) String() string {
	return fmt.Sprintf("server.%d=%s:2888:3888:participant;0.0.0.0:2181", s.id, strings.Split(s.addr, ":")[0])
}

func addServer(conn *zk.Conn, server *Node) error {
	servers := make([]string, 1)
	servers[0] = server.String()
	fmt.Println(servers[0])
	_, err := conn.IncrementalReconfig(servers, nil, -1)
	return err
}

func removeServer(conn *zk.Conn, server *Node) error {
	servers := make([]string, 1)
	servers[0] = strconv.Itoa(server.id)
	_, err := conn.IncrementalReconfig(nil, servers, -1)
	return err
}

type DynamicRegister struct {
	oldConn    *zk.Conn
	newConn    *zk.Conn
	oldCluster []*Node
	newCluster []*Node
	entries    []*entry
}

type entry struct {
	conn *zk.Conn
	node *Node
	mode int16
}

func (j *entry) undo() {
	switch j.mode {
	case Remove:
		_ = addServer(j.conn, j.node)
	case Add:
		_ = removeServer(j.conn, j.node)
	default:
		panic("unreachable")
	}
}

func NewDynamicRegister(oldAddr string, newAddr string) (*DynamicRegister, error) {
	oldCluster, oldServers := makeServers(oldAddr, 1)
	newCluster, newServers := makeServers(newAddr, len(oldCluster)+1)

	oc, _, err := zk.Connect(oldServers, time.Second) //*10)
	if err != nil {
		return nil, err
	}
	nc, _, err := zk.Connect(newServers, time.Second) //*10)
	if err != nil {
		return nil, err
	}
	return &DynamicRegister{
		oldConn:    oc,
		newConn:    nc,
		oldCluster: oldCluster,
		newCluster: newCluster,
	}, nil
}

func (m *DynamicRegister) RegisterOldNodesInNewCluster() error {
	err := m.AddServers(m.newConn, m.oldCluster)
	if err != nil {
		return err
	}
	return nil
}

func (m *DynamicRegister) RegisterNewNodesInOldCluster() error {
	err := m.AddServers(m.oldConn, m.newCluster)
	if err != nil {
		return err
	}
	return nil
}

func (m *DynamicRegister) AddServers(conn *zk.Conn, nodes []*Node) error {
	for _, node := range nodes {
		err := addServer(conn, node)
		if err != nil {
			m.recover()
			return err
		}
		m.entries = append(m.entries, &entry{
			conn: conn,
			node: node,
			mode: Add,
		})
	}
	return nil
}

func (m *DynamicRegister) RemoveServers(conn *zk.Conn, nodes []*Node) error {
	for _, node := range nodes {
		err := removeServer(conn, node)
		if err != nil {
			m.recover()
			return err
		}
		m.entries = append(m.entries, &entry{
			conn: conn,
			node: node,
			mode: Remove,
		})
	}
	return nil
}

func (m *DynamicRegister) recover() {
	for i := len(m.entries) - 1; i >= 0; i-- {
		m.entries[i].undo()
	}
}

func makeServers(headless string, cur int) ([]*Node, []string) {
	addrs := strings.Split(headless, ",")
	sort.Strings(addrs)
	nodes := make([]*Node, len(addrs))
	servers := make([]string, len(addrs))
	for i, addr := range addrs {
		addrTrim := strings.TrimSpace(addr)
		nodes[i] = &Node{
			id:   cur,
			addr: addrTrim,
		}
		servers[i] = addrTrim
		cur++
	}
	return nodes, servers
}
