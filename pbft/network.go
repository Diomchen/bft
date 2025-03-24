// 定义网络结构
package pbft

import (
	"fmt"
	"log"
	"sync"
)

// 节点
type Node struct {
	ID         int
	State      State
	View       int                     // 当前视图
	Log        []Message               // 日志
	PrepareMsg map[string]map[int]bool // 准备消息
	CommitMsg  map[string]map[int]bool // 提交消息
	Mutex      sync.Mutex
	Network    *Network
}

type Network struct {
	Nodes      []*Node
	TotalNodes int
	Faulty     int
}

// 新建网络
func NewNetwork(total, faulty int) *Network {
	if faulty > (total-1)/3 {
		log.Fatal("Too many faulty nodes")
	}

	// Create network instance
	network := &Network{
		Nodes:      make([]*Node, total),
		TotalNodes: total,
		Faulty:     faulty,
	}

	// Initialize each node with the network reference
	for i := 0; i < total; i++ {
		network.Nodes[i] = &Node{
			ID:         i,
			State:      NORMAL,
			View:       0,
			Log:        make([]Message, 0),
			PrepareMsg: make(map[string]map[int]bool),
			CommitMsg:  make(map[string]map[int]bool),
			Mutex:      sync.Mutex{},
			Network:    network,
		}

		if i < faulty {
			network.Nodes[i].State = BYZANTINE
			fmt.Printf("Node %d is Byzantine\n", i)
		}
	}

	return network
}
