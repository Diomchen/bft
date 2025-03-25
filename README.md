# 拜占庭容错验证

## 介绍

**拜占庭容错（Byzantine Fault Tolerance，BFT）** 是分布式系统的一种容错机制，它通过将系统分为多个子系统，并使得这些子系统在出现故障时可以继续工作，从而保证系统的可用性。

**实用拜占庭容错（Practical Byzantine Fault Tolerance，PBFT）** 是一种基于拜占庭将军问题的容错算法，它是一种基于消息传递的算法，可以容忍任意数量的恶意节点，并保证系统的最终一致性。（即通过通信换取一致性）

## 原理
![](https://fastly.jsdelivr.net/gh/Diomchen/PiCor/20250325142338.png)

通信过程：
- **REQUEST**：客户端向主节点发送请求 
- **PRE_PREPARE**：主节点广播
- **PREPARE**：节点通信认证，当收到消息数大于 2*f+1 时，节点进入提交阶段

- **COMMIT**: 节点广播，当收到消息数大于 2*f+1 时，节点进入响应阶段
- **REPLY**：响应消息，消息落盘

## 实现

**v0.1**
该版本是基于以下背景进行模拟的：
- 主节点默认为可信节点，即没有 view change 过程
- 系统中恶意节点不发送信息

![](https://fastly.jsdelivr.net/gh/Diomchen/PiCor/20250325141744.png)

注意：
1. 计算 2*f+1 消息数时，需要将自己的消息也加入。
2. 该版本模拟中主节点需要先处理自己的请求，然后再广播消息，因为没有设置额外的客户端。


## 执行
```go
go run main.go
```