package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

func main() {
	c, _, err := zk.Connect([]string{"10.109.0.15"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
	children, stat, ch, err := c.ChildrenW("/")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v %+v\n", children, stat)
	e := <-ch
	fmt.Printf("%+v\n", e)
	servers := make([]string, 0)
	servers = append(servers, "server.4=zk-cluster-8nxyny-3.zk-cluster-8nxyny-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2888:3888:participant;0.0.0.0:2181")
	reconfig, err := c.Reconfig(servers, -1)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(reconfig)
}
