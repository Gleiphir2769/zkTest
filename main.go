package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"time"
)

func main() {
	c, _, err := zk.Connect([]string{"10.109.58.101"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
	//children, stat, ch, err := c.ChildrenW("/")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%+v %+v\n", children, stat)
	//e := <-ch
	//fmt.Printf("%+v\n", e)
	server := &Server{
		id:   1,
		addr: "zk-cluster-5a914f-0.zk-cluster-5a914f-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2888:3888:participant;0.0.0.0:2181",
	}

	err = AddServer(c, server)
	fmt.Println("-----------------")
	if err != nil {
		fmt.Println("operation success")
	}
}

type Server struct {
	id   int
	addr string
}

func (s Server) String() string {
	return fmt.Sprintf("server.%d=%s", s.id, s.addr)
}

func AddServer(conn *zk.Conn, server *Server) error {
	servers := make([]string, 1)
	servers[0] = server.String()
	_, err := conn.IncrementalReconfig(servers, nil, -1)
	return err
}

func RemoveServer(conn *zk.Conn, server *Server) error {
	servers := make([]string, 1)
	servers[0] = strconv.Itoa(server.id)
	_, err := conn.IncrementalReconfig(nil, servers, -1)
	return err
}
