package pbft

// 消息类型
type MessageType int

const (
	REQUEST MessageType = iota
	PRE_PREPARE
	PREPARE
	COMMIT
	REPLY
)

// 节点状态
type State int

const (
	NORMAL State = iota
	BYZANTINE
)

// 消息结构
type Message struct {
	Type      MessageType
	From      int
	To        int
	Timestamp int64
	Content   string
	Digest    string
}
