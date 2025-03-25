package pbft

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// 创建新节点
func NewNode(id int, state State, network *Network) *Node {
	return &Node{
		ID:         id,
		State:      state,
		View:       0,
		Log:        make([]Message, 0),
		PrepareMsg: make(map[string]map[int]bool),
		CommitMsg:  make(map[string]map[int]bool),
		Network:    network,
	}
}

// 计算消息哈希
func (n *Node) Digest(msg string) string {
	h := sha256.New()
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil))
}

func (n *Node) Broadcast(msg Message, network Network) {
	if n.State == BYZANTINE {
		msg.Content = fmt.Sprintf("Byzantine msg from node: %d", n.ID)
	}

	for _, node := range network.Nodes {
		msgCopy := msg
		msgCopy.To = node.ID
		go node.ReceiveMessage(msgCopy, network)
	}
}

// 计算最大容错节点数f
func F(n int) int {
	return (n - 1) / 3
}

// 处理接收到的消息
func (n *Node) ReceiveMessage(msg Message, network Network) {
	n.Mutex.Lock()
	defer n.Mutex.Unlock()

	switch msg.Type {
	case REQUEST:
		// 这里暂时考虑主节点正常，广播 pre-prepare 消息
		// TODO：后续补充主节点不正常，view change 逻辑
		if n.ID == n.View%len(network.Nodes) && n.State == NORMAL {
			fmt.Printf("【REQUEST】Node %d received request: %s\n", n.ID, msg.Content)
			digest := n.Digest(msg.Content)
			preprepareMsg := Message{
				Type:      PRE_PREPARE,
				From:      n.ID,
				To:        -1,
				Content:   msg.Content,
				Timestamp: time.Now().UnixNano(),
				Digest:    digest,
			}
			n.Broadcast(preprepareMsg, network)
		}

	case PRE_PREPARE:
		if n.State == NORMAL {
			fmt.Printf("【PRE_PREPARE】Node %d received request: %s\n", n.ID, msg.Content)

			digest := n.Digest(msg.Content)
			prepareMsg := Message{
				Type:      PREPARE,
				From:      n.ID,
				To:        -1,
				Content:   msg.Content,
				Timestamp: time.Now().UnixNano(),
				Digest:    digest,
			}
			n.Broadcast(prepareMsg, network)
		}
	case PREPARE:
		if n.State == NORMAL {
			digest := msg.Digest
			if _, exists := n.PrepareMsg[digest]; !exists {
				n.PrepareMsg[digest] = make(map[int]bool)
			}
			n.PrepareMsg[digest][msg.From] = true
			fmt.Printf("【PREPARE】Node %d received request from node %d, has received %d messages\n", n.ID, msg.From, len(n.PrepareMsg[digest]))

			if len(n.PrepareMsg[digest]) >= 2*F(len(network.Nodes))+1 {
				commitMsg := Message{
					Type:      COMMIT,
					From:      n.ID,
					To:        -1,
					Content:   msg.Content,
					Timestamp: time.Now().UnixNano(),
					Digest:    digest,
				}
				fmt.Printf("=> Node %d started to Broadcast COMMIT message\n", n.ID)
				n.Broadcast(commitMsg, network)
			}
		}

	case COMMIT:
		if n.State == NORMAL {
			digest := msg.Digest
			if _, exists := n.CommitMsg[digest]; !exists {
				n.CommitMsg[digest] = make(map[int]bool)
			}
			n.CommitMsg[digest][msg.From] = true

			fmt.Printf("【COMMIT】Node %d received request from node %d, has received %d messages\n", n.ID, msg.From, len(n.CommitMsg[digest]))
			if len(n.CommitMsg[digest]) >= 2*F(len(network.Nodes))+1 {
				replyMsg := Message{
					Type:      REPLY,
					From:      n.ID,
					To:        -1,
					Content:   msg.Content,
					Timestamp: time.Now().UnixNano(),
					Digest:    digest,
				}

				// 直接打印，原理上应该需要返回给客户端
				fmt.Printf("Node %d reached consensus on message: %s\n", n.ID, msg.Content)

				// 发送回复
				if msg.From < len(network.Nodes) {
					network.Nodes[msg.From].ReceiveMessage(replyMsg, network)
				}
			}
		}
	}
}
