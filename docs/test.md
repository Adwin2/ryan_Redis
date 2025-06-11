# 测试

## 测试1

```bash
echo -e "*3\r\n$3\r\nset\r\n$3\r\nfoo\r\n$3\r\n$3\r\nbar\r\n" | nc localhost 6379
```

即 `echo -e "<RESP>" | nc localhost 6379`
