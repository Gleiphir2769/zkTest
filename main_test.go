package main

import (
	"fmt"
	"testing"
)

func TestMake(t *testing.T) {
	nodes := makeServersByHeadless("zk-cluster-z9gphk-0.zk-cluster-z9gphk-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181,zk-cluster-z9gphk-1.zk-cluster-z9gphk-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181,zk-cluster-z9gphk-2.zk-cluster-z9gphk-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2181", 4)
	for _, node := range nodes {
		fmt.Println(node.String())
	}
	// server.1=zk-cluster-z9gphk-0.zk-cluster-z9gphk-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2888:3888:participant:0.0.0.0:2181
	// server.1=zk-cluster-fhdq6s-0.zk-cluster-fhdq6s-headless.zk-cluster.svc.cluster5.nbj04.corp.yodao.com:2888:3888:participant;0.0.0.0:2181
	n_nodes := nodes[2:3]
	for _, node := range n_nodes {
		fmt.Println(node.String())
	}
}
