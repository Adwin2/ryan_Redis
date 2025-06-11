# ryan_Redis 配置指南

> AI-Generated there are some function still on the Todo list.
本文档详细介绍了 ryan_Redis 的配置系统，包括配置文件格式、配置项说明、配置加载优先级以及使用示例。

## 目录

- [配置文件格式](#配置文件格式)
- [配置项说明](#配置项说明)
  - [服务器配置](#服务器配置)
  - [副本配置](#副本配置)
  - [日志配置](#日志配置)
  - [性能调优](#性能调优)
- [配置加载优先级](#配置加载优先级)
  - [环境变量命名规则](#环境变量命名规则)
- [配置示例](#配置示例)
  - [最小化配置](#最小化配置)
  - [完整配置示例](#完整配置示例)
  - [使用环境变量](#使用环境变量)
  - [命令行参数](#命令行参数)
- [最佳实践](#最佳实践)
  - [开发环境](#开发环境)
  - [生产环境](#生产环境)
  - [安全建议](#安全建议)
  - [调试技巧](#调试技巧)
- [故障排除](#故障排除)
- [扩展配置](#扩展配置)
  - [日志配置](#日志配置-1)
  - [性能调优](#性能调优-1)
- [版本兼容性](#版本兼容性)

## 配置文件格式

ryan_Redis 支持 YAML 格式的配置文件，默认配置文件名为 `config.yaml`。配置文件的查找路径如下（按优先级从高到低）：

1. 当前工作目录
2. `./configs/` 目录
3. `/etc/ryan_redis/` 目录
4. `$HOME/.config/ryan_redis/` 目录

## 配置项说明

### 服务器配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `server.port` | string | `6379` | 服务监听端口 |
| `server.bind` | string | `0.0.0.0` | 绑定的IP地址 |
| `server.dir` | string | `/var/lib/rdb` | 持久化数据存储目录 |
| `server.dbfilename` | string | `dump.rdb` | 数据库文件名 |
| `server.role` | string | `master` | 服务器角色：`master` 或 `slave` |
| `server.requirepass` | string | `""` | 认证密码，留空表示不需要认证 |
| `server.maxclients` | int | `10000` | 最大客户端连接数 |
| `server.timeout` | int | `0` | 客户端空闲超时时间（秒），0表示不超时 |

### 副本配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `replica.master_host` | string | `""` | 主节点地址 |
| `replica.master_port` | string | `""` | 主节点端口 |
| `replica.master_auth` | string | `""` | 主节点认证密码 |
| `replica.repl_ping_slave_period` | int | `10` | 从节点ping主节点的间隔（秒） |
| `replica.repl_timeout` | int | `60` | 复制超时时间（秒） |

### 日志配置

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `log.level` | string | `info` | 日志级别：debug, info, warn, error |
| `log.file` | string | `""` | 日志文件路径，空表示输出到标准输出 |
| `log.max_size` | int | `100` | 日志文件最大大小（MB） |
| `log.max_backups` | int | `7` | 保留的旧日志文件数量 |
| `log.max_age` | int | `30` | 保留日志的最大天数 |
| `log.compress` | bool | `true` | 是否压缩旧日志 |

### 性能调优

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `performance.max_memory` | string | `0` | 最大内存使用量，0表示不限制 |
| `performance.max_memory_policy` | string | `noeviction` | 内存淘汰策略 |
| `performance.max_clients` | int | `10000` | 最大客户端连接数 |
| `performance.tcp_keepalive` | int | `300` | TCP keepalive时间（秒） |
| `performance.timeout` | int | `0` | 客户端空闲超时（秒） |
| `performance.tcp_nodelay` | bool | `true` | 是否启用TCP_NODELAY |

## 配置加载优先级

ryan_Redis 的配置加载遵循以下优先级（从高到低）：

1. **命令行参数**：通过 `--flag=value` 形式指定
2. **环境变量**：以 `RR_` 为前缀的环境变量
3. **配置文件**：`config.yaml` 中的配置
4. **默认值**：代码中定义的默认值

### 环境变量命名规则

环境变量名由以下部分组成：
- 前缀：`RR_`
- 配置节：`SERVER_` 或 `REPLICA_` 或 `LOG_` 等
- 配置项：全大写，下划线分隔

例如：
- `RR_SERVER_PORT` 对应 `server.port`
- `RR_REPLICA_MASTER_HOST` 对应 `replica.master_host`
- `RR_LOG_LEVEL` 对应 `log.level`

## 配置示例

### 最小化配置

```yaml
# config.yaml
server:
  port: 6379
  role: master
```

### 完整配置示例

```yaml
# config.yaml
server:
  port: 6379
  bind: 0.0.0.0
  dir: /var/lib/redis
  dbfilename: dump.rdb
  role: slave
  requirepass: "your_secure_password"
  maxclients: 10000
  timeout: 300

replica:
  master_host: 127.0.0.1
  master_port: 6380
  master_auth: "your_master_password"
  repl_ping_slave_period: 10
  repl_timeout: 60

log:
  level: info
  file: /var/log/redis/redis.log
  max_size: 100
  max_backups: 7
  max_age: 30
  compress: true

performance:
  max_memory: "4gb"
  max_memory_policy: "allkeys-lru"
  max_clients: 10000
  tcp_keepalive: 300
  timeout: 0
  tcp_nodelay: true
```

### 使用环境变量

```bash
# 服务器配置
export RR_SERVER_PORT=6380
export RR_SERVER_ROLE=slave
export RR_SERVER_REQUIREPASS=your_secure_password

# 副本配置
export RR_REPLICA_MASTER_HOST=127.0.0.1
export RR_REPLICA_MASTER_PORT=6379
export RR_REPLICA_MASTER_AUTH=your_master_password

# 日志配置
export RR_LOG_LEVEL=debug
export RR_LOG_FILE=/var/log/redis/redis.log
```

### 命令行参数

```bash
# 启动主节点
./ryan_redis \
  --port 6379 \
  --role master \
  --requirepass your_secure_password

# 启动从节点
./ryan_redis \
  --port 6380 \
  --role slave \
  --replicaof 127.0.0.1 6379 \
  --masterauth your_master_password
```

## 最佳实践

### 开发环境

1. 在项目根目录创建 `configs/` 目录
2. 添加 `configs/config.yaml` 文件（添加到 `.gitignore`）
3. 提供 `configs/config.example.yaml` 作为模板

### 生产环境

1. 使用 `/etc/ryan_redis/config.yaml` 作为主配置文件
2. 通过环境变量覆盖敏感配置
3. 使用配置中心管理配置（如 Consul, etcd）

### 安全建议

1. 不要将敏感信息（如密码）直接写入配置文件
2. 使用环境变量或密钥管理服务管理敏感信息
3. 限制配置文件的访问权限：
   ```bash
   chmod 600 /etc/ryan_redis/config.yaml
   chown redis:redis /etc/ryan_redis/config.yaml
   ```

### 调试技巧

1. 查看加载的配置：
   ```bash
   ./ryan_redis --debug
   ```

2. 打印所有配置：
   ```go
   import "github.com/spf13/viper"
   
   func main() {
       viper.Debug()
       // ...
   }
   ```

## 故障排除

### 配置未生效

1. 检查配置文件的查找路径
2. 检查环境变量命名是否正确
3. 使用 `--debug` 标志查看加载的配置

### 配置格式错误

1. 确保 YAML 格式正确
2. 使用在线 YAML 验证器检查语法
3. 检查缩进是否正确

### 权限问题

1. 确保进程有权限读取配置文件
2. 检查数据目录的写入权限
3. 检查日志目录的写入权限

## 扩展配置

### 日志配置

```yaml
log:
  level: info  # debug, info, warn, error
  file: /var/log/redis/redis.log
  max_size: 100  # MB
  max_backups: 7
  max_age: 30    # days
  compress: true  # 是否压缩旧日志
  json: false     # 是否输出JSON格式日志
```

### 性能调优

```yaml
performance:
  max_memory: "4gb"  # 最大内存使用量，支持单位：b, kb, mb, gb
  max_memory_policy: "allkeys-lru"  # 内存淘汰策略
  max_clients: 10000  # 最大客户端连接数
  tcp_keepalive: 300  # TCP keepalive时间（秒）
  timeout: 0  # 客户端空闲超时（秒），0表示不超时
  tcp_nodelay: true  # 是否启用TCP_NODELAY
  hz: 10  # 服务器频率（1-500），值越高CPU使用率越高
```

## 版本兼容性

| ryan_Redis 版本 | 配置文件版本 | 说明 |
|----------------|-------------|------|
| v1.0.0         | v1          | 初始版本 |

> 注意：配置格式在主要版本更新时可能会发生变化，请参考对应版本的文档。