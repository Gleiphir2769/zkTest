package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	oldHeadless := "zk-cluster-fhdq6s-0.zk-cluster-fhdq6s-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181,zk-cluster-fhdq6s-1.zk-cluster-fhdq6s-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181,zk-cluster-fhdq6s-2.zk-cluster-fhdq6s-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181"
	oldInternal := "10.109.12.141"

	newHeadless := "zk-cluster-z9gphk-0.zk-cluster-z9gphk-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181,zk-cluster-z9gphk-1.zk-cluster-z9gphk-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181,zk-cluster-z9gphk-2.zk-cluster-z9gphk-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181"
	newInterNal := "10.109.39.27"

	m, err := NewMigrator(oldHeadless, newHeadless, oldInternal, newInterNal)
	if err != nil {
		panic(err)
	}

	// step1 add old servers to new nodes
	//err = m.AddServers(m.newConn, m.oldCluster[0:2])
	//if err != nil {
	//	panic(err)
	//}

	// step2 redeploy

	err = m.AddServers(m.newConn, m.oldCluster[2:3])
	if err != nil {
		panic(err)
	}

	//// step3 add new servers to old nodes
	//err = m.AddServers(m.oldConn, m.newCluster[0:2])
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = m.AddServers(m.oldConn, m.newCluster[2:3])
	//if err != nil {
	//	panic(err)
	//}

	// step4 redeploy old cluster followers

	// step5 redeploy old cluster leader

	fmt.Println("-----------------")
	if err != nil {
		fmt.Println("operation failed", err)
	}
	fmt.Println("operation success")
}

type Node struct {
	id   int
	addr string
}

func (s *Node) String() string {
	return fmt.Sprintf("server.%d=%s:2888:3888:participant;0.0.0.0:2181", s.id, s.addr)
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

type Migrator struct {
	oldConn    *zk.Conn
	newConn    *zk.Conn
	oldCluster []*Node
	newCluster []*Node
}

func NewMigrator(oldHeadless string, newHeadless string, oldAddr string, newAddr string) (*Migrator, error) {
	oc, _, err := zk.Connect([]string{oldAddr}, time.Second) //*10)
	if err != nil {
		return nil, err
	}
	nc, _, err := zk.Connect([]string{newAddr}, time.Second) //*10)
	if err != nil {
		return nil, err
	}
	return &Migrator{
		oldConn:    oc,
		newConn:    nc,
		oldCluster: makeServersByHeadless(oldHeadless, 1),
		newCluster: makeServersByHeadless(newHeadless, len(oldHeadless)+1),
	}, nil
}

func (m *Migrator) Migrate() {

}

func (m Migrator) AddServers(conn *zk.Conn, nodes []*Node) error {
	for _, node := range nodes {
		err := addServer(conn, node)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m Migrator) RemoveServers(conn *zk.Conn, nodes []*Node) error {
	for _, node := range nodes {
		err := removeServer(conn, node)
		if err != nil {
			return err
		}
	}
	return nil
}

func makeServersByHeadless(headless string, cur int) []*Node {
	addrs := strings.Split(headless, ",")
	sort.Strings(addrs)
	nodes := make([]*Node, len(addrs))
	for i, addr := range addrs {
		addrTrim := strings.TrimSpace(addr)
		nodes[i] = &Node{
			id:   cur,
			addr: strings.Split(addrTrim, ":")[0],
		}
		cur++
	}
	return nodes
}
