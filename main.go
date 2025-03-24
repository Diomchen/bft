package main

import (
	"bft_valid/pbft"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

func main() {
	// 初始化网络
	totalNodes := 10
	byzantineNodes := 2
	network := pbft.NewNetwork(totalNodes, byzantineNodes)

	// 选择客户端节点(这里不考虑客户端为 Byzantine 的情况)
	clientNode := 3

	// 创建 REQUEST 消息
	request := pbft.Message{
		Type:      pbft.REQUEST,
		From:      clientNode,
		To:        -1,
		Timestamp: time.Now().UnixNano(),
		Content:   "Transfer 100 tokens from A to B.",
		Digest:    "",
	}

	// 计算摘要
	h := sha256.New()
	h.Write([]byte(request.Content))
	request.Digest = hex.EncodeToString(h.Sum(nil))

	// 广播 REQUEST 消息
	fmt.Printf("Broadcasting REQUEST message from node %d to node %d\n", clientNode, request.To)
	network.Nodes[clientNode].Broadcast(request, *network)

	// 模拟共识等待过程
	time.Sleep(3 * time.Second)

	// 验证结果
	fmt.Println("\n--- Verification Results ---")
	consensus := make(map[string]int)

	for i, node := range network.Nodes {
		if node.State == pbft.NORMAL {
			for digest, commits := range node.CommitMsg {
				if len(commits) >= 2*pbft.F(totalNodes)+1 {
					consensus[digest]++
					fmt.Printf("Node %d reached consensus on digest: %s\n", i, digest)
				}
			}
		}
	}

	fmt.Println("\n--- Final Consensus ---")
	for digest, count := range consensus {
		fmt.Printf("Digest %s: %d normal nodes reached consensus\n", digest, count)
	}
}
