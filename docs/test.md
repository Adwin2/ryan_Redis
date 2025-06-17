# 测试

## 编译并启动主从节点

```bash # 1
go build -o ryan_redis ./main.go
```

```bash # 2
./ryan_redis
```

```bash # 3
./ryan_redis -p 6380 --replicaof "0.0.0.0 6379"
```

## 测试SET命令

```bash # 1
echo -e "*3\r\n$3\r\nset\r\n$3\r\nfoo\r\n$3\r\n$3\r\nbar\r\n" | nc localhost 6379
```

> 即 `echo -e "<RESP>" | nc localhost 6379`

## 测试GET命令

```bash # 1
echo -e "*2\r\n$3\r\nget\r\n$3\r\nfoo\r\n" | nc localhost 6379
```
