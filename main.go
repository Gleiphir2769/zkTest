package main

import "fmt"

//const (
//	Remove = 1
//	Add    = 2
//)

func main() {
	oldHeadless := "zk-cluster-b2vt9x-0.zk-cluster-b2vt9x-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181,zk-cluster-b2vt9x-2.zk-cluster-b2vt9x-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181,zk-cluster-b2vt9x-1.zk-cluster-b2vt9x-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181"

	newHeadless := "zk-cluster-xcf8le-1.zk-cluster-xcf8le-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181,zk-cluster-xcf8le-2.zk-cluster-xcf8le-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181,zk-cluster-xcf8le-0.zk-cluster-xcf8le-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181,zk-cluster-xcf8le-4.zk-cluster-xcf8le-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181,zk-cluster-xcf8le-3.zk-cluster-xcf8le-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181"

	m, err := NewDynamicRegister(oldHeadless, newHeadless)
	if err != nil {
		panic(err)
	}

	err = m.RegisterOldNodesInNewCluster()
	fmt.Println("-----------------")
	if err != nil {
		fmt.Println("operation failed", err)
	}
	fmt.Println("operation success")

	//// step1 add first and second old servers to new nodes
	//err = m.AddServers(m.newConn, m.oldCluster[0:2])
	//if err != nil {
	//	panic(err)
	//}
	//
	//// step2 all nodes redeploy
	//
	//// step3 add the last old servers
	//err = m.AddServers(m.newConn, m.oldCluster[2:3])
	//if err != nil {
	//	panic(err)
	//}
	//
	////// step3 add new servers to old nodes
	//err = m.AddServers(m.oldConn, m.newCluster[0:3])
	//fmt.Println("-----------------")
	//if err != nil {
	//	fmt.Println("operation failed", err)
	//}
	//fmt.Println("operation success")
	//
	//// step4 redeploy old cluster followers
	//
	//// step5 redeploy old cluster leader
	//
	//fmt.Println("-----------------")
	//if err != nil {
	//	fmt.Println("operation failed", err)
	//}
	//fmt.Println("operation success")
}

//type Node struct {
//	id   int
//	addr string
//}
//
//func (s *Node) String() string {
//	return fmt.Sprintf("server.%d=%s:2888:3888:participant;0.0.0.0:2181", s.id, s.addr)
//}
//
//func addServer(conn *zk.Conn, server *Node) error {
//	servers := make([]string, 1)
//	servers[0] = server.String()
//	fmt.Println(servers[0])
//	_, err := conn.IncrementalReconfig(servers, nil, -1)
//	return err
//}
//
//func removeServer(conn *zk.Conn, server *Node) error {
//	servers := make([]string, 1)
//	servers[0] = strconv.Itoa(server.id)
//	_, err := conn.IncrementalReconfig(nil, servers, -1)
//	return err
//}
//
//type Migrator struct {
//	oldConn    *zk.Conn
//	newConn    *zk.Conn
//	oldCluster []*Node
//	newCluster []*Node
//	entries    []*entry
//}
//
//type entry struct {
//	conn *zk.Conn
//	node *Node
//	mode int16
//}
//
//func (j *entry) undo() {
//	switch j.mode {
//	case Remove:
//		_ = addServer(j.conn, j.node)
//	case Add:
//		_ = removeServer(j.conn, j.node)
//	default:
//		panic("unreachable")
//	}
//}
//
//func NewMigrator(oldHeadless string, newHeadless string, oldAddr string, newAddr string) (*Migrator, error) {
//	oc, _, err := zk.Connect([]string{oldAddr}, time.Second) //*10)
//	if err != nil {
//		return nil, err
//	}
//	nc, _, err := zk.Connect([]string{newAddr}, time.Second) //*10)
//	if err != nil {
//		return nil, err
//	}
//	return &Migrator{
//		oldConn:    oc,
//		newConn:    nc,
//		oldCluster: makeServersByHeadless(oldHeadless, 1),
//		newCluster: makeServersByHeadless(newHeadless, len(oldHeadless)+1),
//	}, nil
//}
//
//func (m *Migrator) Migrate() {
//	err := m.ModifyNewCluster()
//	if err != nil {
//		log.Fatalln(err)
//	}
//	fmt.Println("------------------")
//	fmt.Println("operation success")
//}
//
//func (m Migrator) ModifyNewCluster() error {
//	//err := m.AddServers(m.newConn, m.oldCluster[0:2])
//	//if err != nil {
//	//	return err
//	//}
//
//	err := m.AddServers(m.newConn, m.oldCluster[2:3])
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (m *Migrator) AddServers(conn *zk.Conn, nodes []*Node) error {
//	for _, node := range nodes {
//		err := addServer(conn, node)
//		if err != nil {
//			m.recover()
//			return err
//		}
//		m.entries = append(m.entries, &entry{
//			conn: conn,
//			node: node,
//			mode: Add,
//		})
//	}
//	return nil
//}
//
//func (m *Migrator) RemoveServers(conn *zk.Conn, nodes []*Node) error {
//	for _, node := range nodes {
//		err := removeServer(conn, node)
//		if err != nil {
//			m.recover()
//			return err
//		}
//		m.entries = append(m.entries, &entry{
//			conn: conn,
//			node: node,
//			mode: Remove,
//		})
//	}
//	return nil
//}
//
//func (m *Migrator) recover() {
//	for i := len(m.entries) - 1; i >= 0; i-- {
//		m.entries[i].undo()
//	}
//}
//
//func makeServersByHeadless(headless string, cur int) []*Node {
//	addrs := strings.Split(headless, ",")
//	sort.Strings(addrs)
//	nodes := make([]*Node, len(addrs))
//	for i, addr := range addrs {
//		addrTrim := strings.TrimSpace(addr)
//		nodes[i] = &Node{
//			id:   cur,
//			addr: strings.Split(addrTrim, ":")[0],
//		}
//		cur++
//	}
//	return nodes
//}
