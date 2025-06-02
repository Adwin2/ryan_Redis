package protocol

import "net"

// ResponseWriter 定义了生成 Redis 协议响应的接口
type ResponseWriter interface {
	// 返回对应连接
	Conn() net.Conn
	// 写入简单字符串响应
	WriteSimpleString(str string) error
	// 写入批量字符串响应
	WriteBulkString(str string) error
	// 写入数组响应
	WriteArray(str []string) error
	// 写入空批量字符串 (-1)
	WriteNull() error
	// 刷新缓冲区
	Flush() error
}
