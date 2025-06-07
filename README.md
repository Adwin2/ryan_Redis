# ryan_Redis

[![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

一个使用 Go 语言实现的简化版 Redis 服务器，支持基础键值存储和主从复制功能。本项目旨在深入理解 Redis 的核心协议、网络通信模型以及数据持久化机制。

## 🚀 功能特性

### 核心功能

- [x] 支持 RESP (REdis Serialization Protocol) 协议
- [x] 基础数据结构支持：String, List, Hash, Set, Sorted Set
- [x] 主从复制（Master-Slave Replication）
- [x] 支持 AOF 持久化
- [ ] 事务支持（开发中）
- [ ] 集群模式（规划中）

### 技术亮点

- 使用 Go 原生 net 包实现高性能网络通信
- 基于事件循环的高并发模型
- 无锁数据结构设计，提高并发性能
- 支持断点续传的主从同步机制

## 📦 快速开始

### 环境要求

- Go 1.20 或更高版本
- Git

### 安装与运行

```bash
# 克隆项目
git clone https://github.com/Adwin2/ryan_Redis.git
cd ryan_Redis

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

### 项目结构

```bash
ryan_Redis/
    └── app
        ├── cmd  # 项目的命令行工具
        |    ├── main.go # 主程序入口
        |    └── server  # 服务端
        ├── internal  # 项目的核心实现
        |    ├── command  # 命令处理器
        |    ├── config  # 配置相关
        |    ├── protocol  # Redis 协议实现
        |    ├── replication  # 主从复制相关
        |    ├── rtest  # 测试帮助函数
        |    ├── server  # 服务端
        |    |    ├── master  # 主节点
        |    |    └── slave  # 从节点
        |    └── storage  # 存储相关
        |         ├── memory  # 内存存储
        |         |    └── kvstore  # KV 存储
        |         └── rdb  # RDB 存储
        └── pkg  # 公共库
            └── errors_r  # 错误处理
```

### 📚 实现细节

1. 网络层
使用 Go 的 net 包实现 TCP 服务器
基于 goroutine 的并发模型
连接池管理

2. 主从复制
支持全量同步和增量同步
断点续传
心跳检测

3. 持久化
AOF 持久化
支持 AOF 重写

### 🧪 测试

```bash
# 运行单元测试
go test -v ./...

# 运行性能测试
go test -bench=. -benchmem
```

### 🤝 贡献

欢迎提交 Issue 和 PR。对于重大更改，请先开启 Issue 讨论您希望更改的内容。

### 📜 开源协议

本项目采用 MIT 许可证 - 详情请参阅 LICENSE 文件。
