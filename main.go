package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
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
	servers := make([]string, 0)
	servers = append(servers, "1")

	reconfig, err := c.IncrementalReconfig(nil, servers, -1)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("-----------------")
	fmt.Println(reconfig)
}
