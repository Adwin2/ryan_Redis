package protocol

import (
	"bufio"
	"net"
)

// connResponseWriter 是基于 net.Conn 的 ResponseWriter 实现
type connResponseWriter struct {
	conn    net.Conn
	writer  *bufio.Writer
	scratch [64]byte // 用于小量数据的临时缓冲区
}

// NewConnResponseWriter 创建一个新的基于连接的 ResponseWriter
func NewConnResponseWriter(conn net.Conn) ResponseWriter {
	return &connResponseWriter{
		conn:   conn,
		writer: bufio.NewWriter(conn),
	}
}

func (w *connResponseWriter) Conn() net.Conn {
	return w.conn
}

// 封装
func (w *connResponseWriter) WriteSimpleString(str string) error {
	_, err := w.conn.Write(SimpleStringFmt(str))
	return err
}

func (w *connResponseWriter) WriteBulkString(str string) error {
	_, err := w.conn.Write(BulkStringFmt(str))
	return err
}

func (w *connResponseWriter) WriteArray(str []string) error {
	_, err := w.conn.Write(ArrayFmt(str))
	return err
}

func (w *connResponseWriter) WriteNull() error {
	_, err := w.conn.Write(NullFmt())
	return err
}

func (w *connResponseWriter) Flush() error {
	return w.writer.Flush()
}
