# ryan_Redis

[![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Documentation](https://img.shields.io/badge/docs-configuration-brightgreen)](docs/configuration.md)

一个使用 Go 语言实现的简化版 Redis 服务器，支持基础键值存储、主从复制功能以及RDB持久化。本项目旨在深入理解 Redis 的核心协议、网络通信模型以及数据持久化机制。

## 功能特性

### 核心功能

- [x] 支持 RESP (REdis Serialization Protocol) 协议
- [x] 数据结构支持：String
- [x] 主从复制（Master-Slave Replication）
- [x] 支持 RDB 持久化

### 运行

```bash
# 启动主节点
go run cmd/server/main.go -port 6379 -role master

# 启动从节点（在另一个终端）
go run cmd/server/main.go -port 6380 -role slave -master-addr localhost:6379
```

### 使用示例

```bash
# 连接到主节点
redis-cli -p 6379

# 设置键值
SET mykey "Hello Redis"

# 获取键值
GET mykey

# 查看主从复制状态
INFO replication
```

测试示例 见 ./docs/test.md

### 实现细节

1. 网络层
使用 Go 的 net 包实现 TCP 服务器
基于 goroutine 的并发模型

2. 主从复制
支持全量同步

3. 持久化
RDB 持久化
本项目采用 MIT 许可证- 详情请参阅 [LICENSE](https://opensource.org/licenses/MIT) 文件。
