package main

import (
	"fmt"
	"testing"
)

func TestMake(t *testing.T) {
	nodes := makeServersByHeadless("zk-cluster-z9gphk-1.zk-cluster.cluster5.nbj04.corp.yodao.com:2181,zk-cluster-z9gphk-2.zk-cluster.cluster5.nbj04.corp.yodao.com:2181,zk-cluster-z9gphk-0.zk-cluster.cluster5.nbj04.corp.yodao.com:2181", 4)
	for _, node := range nodes {
		fmt.Println(node.addr)
	}

	n_nodes := nodes[2:3]
	for _, node := range n_nodes {
		fmt.Println(node.addr)
	}
}
